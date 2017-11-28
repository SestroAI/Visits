package endpoints

import (
	"github.com/emicklei/go-restful"

	"errors"
	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/dao"
	"github.com/SestroAI/shared/logger"
	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/models/visits"
	"github.com/SestroAI/shared/routing"
	serrors "github.com/SestroAI/shared/utils/errors"
	"net/http"
	"github.com/SestroAI/Visits/controller/events"
)

type VisitResource struct {
}

func (u VisitResource) Register(container *restful.Container, prefix string) {
	ws := new(restful.WebService)

	ws.Path(prefix + "/visits").
		Doc("Manage Sestro User Visits").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		ApiVersion("1.0.0")

	ws.Route(
		ws.POST("/").To(u.CreateVisit).
			Doc("Start a new Visit").
			Operation("CreateVisit").
			Reads(NewVisitData{}).
			Writes(NewVisitResponse{}).
			Filter(routing.IdempotencyFilter))

	ws.Route(
		ws.GET("/{visit_id}").To(u.GetVisit).
			Doc("Get Visit Info").
			Operation("GetVisit").
			Param(ws.PathParameter("visit_id", "Merchant Visit ID").DataType("string")).
			Writes(visits.MerchantVisit{}))

	ws.Route(
		ws.POST("/{visit_id}/end").To(u.EndVisit).
			Param(ws.PathParameter("visit_id", "Merchant Visit ID").DataType("string")).
			Reads(EndVisitInput{}))
	container.Add(ws)
}

func (u VisitResource) GetVisit(req *restful.Request, res *restful.Response) {
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)
	visit, err := ref.GetVisit(req.PathParameter("visit_id"))
	if err != nil {
		logger.ReqInfof(req, "Unable to get visit with id = %s and error = %s",
			req.PathParameter("visit_id"), err.Error())
		res.WriteErrorString(http.StatusNotFound, "Unable to get visit data err = "+err.Error())
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, visit)
}

type NewVisitData struct {
	RestaurantId string `json:"merchantId"`
	TableId      string `json:"tableId"`
	GeoLatitude  string `json:"geoLatitude"`
	GeoLongitude string `json:"geoLongitude"`
}

func (visit_data *NewVisitData) Verify() error {
	if visit_data.RestaurantId == "" || visit_data.TableId == "" {
		return errors.New("No tableId or restaurantId found!")
	}
	return nil
}

type NewVisitResponse struct {
	VisitId        string `json:"visitId"`
	DinerSessionId string `json:"sessionId"`
}

/*
	Returns an ongoing visit for the table if exists,
	throws conflict if diner is already a part of another visit currently,
	checks for fradulent behaviour
	creates a new visit if none of above is true
*/
func (u VisitResource) CreateVisit(req *restful.Request, res *restful.Response) {
	diner, ok := req.Attribute(config.RequestUser).(*auth.User)
	if !ok || diner.ID == "" {
		logger.ReqErrorf(req, "No Diner found")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}
	logger.Infof("Diner Id = %s", diner.ID)
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)

	var data NewVisitData
	err := req.ReadEntity(&data)
	if err != nil {
		logger.ReqInfof(req, "Incompatible data type for New Visit Data. Unable to create Visit.")
		res.WriteError(http.StatusBadRequest, serrors.ErrWrongDataFormat)
		return
	}

	//Verify that diner is in the Restaurant by location
	err = data.Verify()
	if err != nil {
		logger.ReqErrorf(req, "Fraudulent Request Found. Unable to verify request data")
		res.WriteError(http.StatusBadRequest, err)
		return
	}

	//Check if any existing visit in progress for requesting diner
	if diner.CustomerProfile.OngoingVisitId != "" {
		//Found. Return the existing visit
		res.WriteErrorString(http.StatusConflict, "An active visit is already going on for this user. "+
			"Please end that to start a new one.")
		return
	}

	//Since visit does not already exists, start a new one
	visit, err := ref.StartNewVisit(diner, data.TableId)
	if err != nil {
		if visit != nil {
			res.WriteHeader(http.StatusAlreadyReported)
		} else {
			if err == serrors.ErrConflict {
				res.WriteError(http.StatusBadRequest, serrors.ErrBadData)
			} else {
				logger.ReqErrorf(req, "Unable to start a new Visit for table: %s and err = %s",
					data.TableId, err.Error())
				res.WriteError(http.StatusInternalServerError, serrors.ErrServerError)
			}
			return
		}
	} else {
		res.WriteHeader(http.StatusOK)
	}

	dinerSessionId, ok := visit.Diners[diner.ID]
	if !ok {
		//This should never happen
		logger.ReqErrorf(req, "Unable to get Diner (id = %s) from new visit object (id = %s)", diner.ID, visit.ID)
		res.WriteError(http.StatusInternalServerError, serrors.ErrServerError)
		return
	}

	res.WriteEntity(
		NewVisitResponse{
			VisitId:        visit.ID,
			DinerSessionId: dinerSessionId,
		})
	return
}

type EndVisitInput struct {
	visits.Rating `json:"guestRating"`
}

func (u VisitResource) EndVisit(req *restful.Request, res *restful.Response) {
	token, _ := req.Attribute(config.RequestToken).(string)
	user, _ := req.Attribute(config.RequestUser).(*auth.User)

	ref := dao.NewVisitDao(token)

	input := EndVisitInput{}
	err := req.ReadEntity(&input)
	if err != nil {
		res.WriteErrorString(http.StatusBadRequest, "Wrong data input")
		return
	}

	visitId := req.PathParameter("visit_id")

	visit, err := ref.GetVisit(visitId)
	if err != nil {
		res.WriteErrorString(http.StatusInternalServerError, "Unable to get visit ID = " + visitId)
		return
	}

	if visit.IsComplete {
		res.WriteHeaderAndEntity(http.StatusAlreadyReported, "Visit has already been ended")
		return
	}

	visit.GuestRating = &input.Rating
	err = ref.EndVisit(visit)
	if err != nil {
		logger.ReqErrorf(req, "unable to end visit ID = %s with error = %s", visit.ID, err.Error())
		res.WriteErrorString(http.StatusInternalServerError, "Cannot end the visit right now")
		return
	}

	err = events.SendEndVisitEvent(user.ID, visit)
	if err != nil {
		logger.ReqErrorf(req, "Unable to send visit end for visit ID = %s message with err = %s",
			visit.ID, err.Error())
		/*
		TODO: Dont know what to do here!
		 */
		 res.WriteErrorString(http.StatusInternalServerError, "Unable to send visit completion event")
		 return
	}

	res.WriteHeader(http.StatusOK)
	return
}