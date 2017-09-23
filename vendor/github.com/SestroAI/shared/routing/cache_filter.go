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

	iid := req.HeaderParameter("Idempotency-Key")
	if iid == "" {
		res.WriteErrorString(http.StatusExpectationFailed, "No Idempotency Id found. Please set the header " +
			"'Idempotency-Key' as a unique ID to enable idempotency. You can use this id to make duplicate calls and " +
			"all requests will be idempotent")
		return
	}
	sr_conn := sestro_redis.GetNewRedisConnection()

	redis_key := ""


	user, ok := req.Attribute(config.RequestUser).(*auth.User)
	if ok && user != nil {
		redis_key = user.ID + "::" + iid
	} else {
		redis_key = iid
	}

	value, err := sr_conn.GetKeyValueFromRedis(redis_key)
	if err == nil && value != nil {
			int_value, ok := value.(uint8)
			if  ok {
				//This is a repeated IID for this user.
				if int_value == 1 {
					//Last execution was successful, return alreadyreported
					res.WriteHeader(http.StatusAlreadyReported)
					return
				}
			}
	} else {
		//Save iid in cache
		err = sr_conn.SaveKeyValueInRedis(redis_key, 0) //"0" means failure, by default
		if err != nil {
			logger.ReqErrorf(req, "Unable to write in cache error: %s", err.Error())
		}
	}

	//Execute Request
	chain.ProcessFilter(req, res)

	//Post Processing
	if res.Error() == nil {
		//Request was not successful
		err = sr_conn.SaveKeyValueInRedis(redis_key, 1) //"1" means successful
		if err != nil {
			logger.ReqErrorf(req, "Unable to write in cache error: %s", err.Error())
		}
	}
}
