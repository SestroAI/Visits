package endpoints

import (
	"github.com/emicklei/go-restful"

	"github.com/SestroAI/shared/config"
	"github.com/SestroAI/shared/dao"
	"github.com/SestroAI/shared/logger"
	"github.com/SestroAI/shared/models/visits"

	"net/http"
	"github.com/SestroAI/shared/models/orders"
	serrors "github.com/SestroAI/shared/utils/errors"
)

type SessionResource struct {
}

func (u SessionResource) Register(container *restful.Container, prefix string) {
	ws := new(restful.WebService)

	ws.Path(prefix + "/sessions").
		Doc("Manage Sestro User Sessions").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		ApiVersion("1.0.0")

	ws.Route(
		ws.GET("/{session_id}").To(u.GetSession).
			Doc("Get Session Info").
			Operation("GetSession").
			Param(ws.PathParameter("session_id", "User Session ID").DataType("string")).
			Writes(visits.VisitDinerSession{}))

	ws.Route(
		ws.PUT("/{session_id}/orders").To(u.AddOrder).
			Doc("Add items to order").
			Param(ws.PathParameter("session_id", "User Session ID").DataType("string")).
			Operation("AddOrder").
			Reads(AddOrderInput{}).
			Writes(visits.VisitDinerSession{}))

	ws.Route(
		ws.PUT("/{session_id}/order/{order_id}/status/{status}").To(u.UpdateOrderStatus).
			Doc("Update Order Status").
			Param(ws.PathParameter("session_id", "User Session ID").DataType("string")).
			Param(ws.PathParameter("order_id", "Order ID").DataType("string")).
			Writes(visits.VisitDinerSession{}))

	container.Add(ws)
}

func (u SessionResource) GetSession(req *restful.Request, res *restful.Response) {
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)
	session, err := ref.GetVisitSession(req.PathParameter("session_id"))
	if err != nil {
		logger.ReqInfof(req, "Unable to get session with id = %s and error = %s",
			req.PathParameter("session_id"), err.Error())
		res.WriteErrorString(http.StatusNotFound, "Unable to get session data err = "+err.Error())
		return
	}

	res.WriteHeaderAndEntity(http.StatusOK, session)
}

type AddOrderInput struct {
	Items []string `json:"items"`
}

func (u SessionResource) AddOrder(req *restful.Request, res *restful.Response) {
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)
	sessionId := req.PathParameter("session_id")

	session, err := ref.GetVisitSession(sessionId)
	if err != nil {
		res.WriteErrorString(http.StatusNotFound, "No Session with ID = " + sessionId + " found")
		return
	}

	var data AddOrderInput
	err = req.ReadEntity(&data)
	if err != nil {
		logger.ReqErrorf(req, "Invalid data format")
		res.WriteError(http.StatusBadRequest, serrors.ErrWrongDataFormat)
		return
	}

	for _, item := range data.Items {
		/*
		TODO: Check if item is valid for this restaurant and session and available
		 */
		order := orders.NewOrder()
		order.ItemId = item
		//Order status is default which is "delivered"
		session.Orders[order.ID] = order
	}

	err = ref.SaveVisitSession(session.ID, session)
	if err != nil {
		logger.ReqErrorf(req, "Unable to save session ID = %s after adding order with error = %s",
			session.ID, err.Error())
		res.WriteErrorString(http.StatusInternalServerError, "Unable to save the session")
		return
	}
	res.WriteHeader(http.StatusOK)
	res.WriteEntity(session)
	return
}

func (u SessionResource) UpdateOrderStatus(req *restful.Request, res *restful.Response) {
	token, _ := req.Attribute(config.RequestToken).(string)

	ref := dao.NewVisitDao(token)
	sessionId := req.PathParameter("session_id")
	orderId := req.PathParameter("order_id")

	validStatus := false
	for _, possibleStatus := range orders.AllowedOrderStatus {
		if possibleStatus == req.PathParameter("status"){
			validStatus = true
		}
	}
	if !validStatus {
		res.WriteErrorString(http.StatusBadRequest, "Invalid status provided")
		return
	}

	session, err := ref.GetVisitSession(sessionId)
	if err != nil {
		res.WriteErrorString(http.StatusNotFound, "No Session with ID = " + sessionId + " found")
		return
	}

	var currOrder *orders.Order
	found := false
	for _, order := range session.Orders {
		if order.ID == orderId {
			found = true
			currOrder = order
			break
		}
	}
	if !found {
		res.WriteErrorString(http.StatusNotFound, "No Order with ID = " + orderId + " found")
		return
	}

	currOrder.Status = req.PathParameter("status")
	session.Orders[currOrder.ID] = currOrder

	err = ref.SaveVisitSession(session.ID, session)
	if err != nil {
		logger.ReqErrorf(req, "Unable to save session ID = %s after updating order with error = %s",
			session.ID, err.Error())
		res.WriteErrorString(http.StatusInternalServerError, "Unable to save the session")
		return
	}
	res.WriteHeader(http.StatusOK)
	res.WriteEntity(session)
	return
}