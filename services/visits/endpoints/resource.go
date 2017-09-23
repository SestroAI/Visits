package endpoints

import(
	"github.com/emicklei/go-restful"

	"github.com/SestroAI/shared/models/auth"
	"github.com/SestroAI/shared/logger"
	"net/http"
	"github.com/SestroAI/shared/routing"
	"github.com/SestroAI/shared/dao"
	serrors "github.com/SestroAI/shared/utils/errors"
	"errors"
	"github.com/SestroAI/shared/config"
)

type VisitResource struct {
}

func (u VisitResource) Register(container *restful.Container, prefix string)  {
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
			Writes(NewVisitResponse{})).
			Filter(routing.IdempotencyFilter)

	container.Add(ws)
}

type NewVisitData struct {
	RestaurantId string `json:"merchantId"`
	TableId string `json:"tableId"`
	GeoLatitude string `json:"geoLatitude"`
	GeoLongitude string `json:"geoLongitude"`
}

func (visit_data *NewVisitData) Verify() error {
	if visit_data.RestaurantId == "" || visit_data.TableId == ""{
		return errors.New("No tableId or restaurantId found!")
	}
	return nil
}

type NewVisitResponse struct {
	VisitId string `json:"visitId"`
	DinerSessionId string `json:"sessionId"`
}


/*
	Returns an ongoing visit for the table if exists,
	throws conflict if diner is already a part of another visit currently,
	checks for fradulent behaviour
	creates a new visit if none of above is true
*/
func (u VisitResource) CreateVisit(req *restful.Request, res *restful.Response)  {
	diner, _ := req.Attribute(config.RequestDiner).(*auth.User)
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)

	var data NewVisitData
	err := req.ReadEntity(&data)
	if err != nil {
		logger.ReqInfof(req,"Incompatible data type for New Visit Data. Unable to create Visit.")
		res.WriteError(http.StatusBadRequest, serrors.ErrWrongDataFormat)
		return
	}

	//Verify that diner is in the Restaurant by location
	err = data.Verify(); if err != nil {
		logger.ReqErrorf(req, "Fraudulent Request Found. Unable to verify request data")
		res.WriteError(http.StatusBadRequest, err)
		return
	}

	//Check if any existing visit in progress for requesting diner
	if diner.CustomerProfile.OngoingVisitId != "" {
		//Found. Return the existing visit
		res.WriteErrorString(http.StatusConflict,"An active visit is already going on for this user. " +
			"Please end that to start a new one.")
		return
	}

	//Since visit does not already exists, start a new one
	visit, err := ref.StartNewVisit(diner, data.TableId)
	if err != nil {
		if visit != nil {
			res.WriteHeader(http.StatusAlreadyReported)
		} else {
			if err == serrors.ErrConflict{
				res.WriteError(http.StatusBadRequest, serrors.ErrBadData)
			} else {
				logger.ReqErrorf(req,"Unable to start a new Visit for table: %s", data.TableId)
				res.WriteError(http.StatusInternalServerError, serrors.ErrServerError)
			}
			return
		}
	} else {
		res.WriteHeader(http.StatusOK)
	}

	dinerSession, ok := visit.Diners[diner.ID]; if !ok{
		//This should never happen
		logger.ReqErrorf(req,"Unable to get Diner (id = %s) from new visit object (id = %s)", diner.ID, visit.ID)
		res.WriteError(http.StatusInternalServerError, serrors.ErrServerError)
		return
	}

	res.WriteEntity(
		NewVisitResponse{
			VisitId:visit.ID,
			DinerSessionId:dinerSession.ID,
		})
	return
}