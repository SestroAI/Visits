package routing

import (
	"github.com/emicklei/go-restful"
	"github.com/SestroAI/shared/logger"
	"net/http"
	"github.com/SestroAI/shared/utils/sestro_redis"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/config"
)

/*
Only enable this if LoggedIn filter is attached before this.
 */
func IdempotencyFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain){
	user, _ := req.Attribute(config.RequestUser).(*auth.User)

	iid := req.HeaderParameter("Idempotency-Key")
	if iid == "" {
		res.WriteErrorString(http.StatusExpectationFailed, "No Idempotency Id found. Please set the header" +
			"IdempotencyId as a unique ID to enable idempotency. You can use this id to make subsequent calls and " +
			"all requests will be idempotent")
		return
	}
	sr_conn := sestro_redis.GetNewRedisConnection()
	redis_key := user.ID + "::" + iid
	value, err := sr_conn.GetKeyValueFromRedis(redis_key)
	if err == nil && value.(string) != ""{
		//This is a repeated IID for this user.
		value = value.(string)
		if value == "1"{
			//Last execution was successful, return alreadyreported
			res.WriteHeader(http.StatusAlreadyReported)
			return
		}
	} else {
		//Save iid in cache
		err = sr_conn.SaveKeyValueInRedis(redis_key, false)
		if err != nil {
			logger.Errorf("Unable to write in cache error: %s", err.Error())
		}
	}
	final_value := true

	//Execute Request
	chain.ProcessFilter(req, res)

	//Post Processing
	if res.Error() != nil {
		//Request was not successful
		final_value = false
	}
	err = sr_conn.SaveKeyValueInRedis(redis_key, final_value)
	if err != nil {
		logger.Errorf("Unable to write in cache error: %s", err.Error())
	}
}
