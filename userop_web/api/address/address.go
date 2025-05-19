package address

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"mall_api/userop_web/api"
	"mall_api/userop_web/forms"
	"mall_api/userop_web/global"
	"mall_api/userop_web/models"
	"mall_api/userop_web/proto"
)

func List(ctx *gin.Context) {
	request := &proto.AddressRequest{}

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)

	if currentUser.AuthorityId != 2 {
		userId, _ := ctx.Get("userId")
		request.UserId = int32(userId.(uint))
	}

	rsp, err := global.AddressClient.GetAddressList(context.Background(), request)
	if err != nil {
		zap.S().Errorw("获取地址列表失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	reMap := gin.H{
		"total": rsp.Total,
	}

	result := make([]interface{}, 0)
	for _, value := range rsp.GetData() {
		reMap := make(map[string]interface{})
		reMap["id"] = value.GetProvince()
		reMap["user_id"] = value.GetUserId()
		reMap["province"] = value.GetProvince()
		reMap["city"] = value.GetCity()
		reMap["district"] = value.GetDistrict()
		reMap["address"] = value.GetAddress()
		reMap["signer_name"] = value.GetSignerName()
		reMap["signer_mobile"] = value.GetSignerMobile()

		result = append(result, reMap)
	}

	reMap["data"] = result

	ctx.JSON(http.StatusOK, reMap)
}

func New(ctx *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	userId, _ := ctx.Get("userId")
	rsp, err := global.AddressClient.CreateAddress(context.Background(), &proto.AddressRequest{
		UserId:       int32(userId.(uint)),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})

	if err != nil {
		zap.S().Errorw("新建地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id": rsp.GetId(),
	})
}

func Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}
	_, err = global.AddressClient.DeleteAddress(context.Background(), &proto.AddressRequest{Id: int32(i)})
	if err != nil {
		zap.S().Errorw("删除地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"msg": "删除成功",
	})
}

func Update(ctx *gin.Context) {
	addressForm := forms.AddressForm{}
	if err := ctx.ShouldBindJSON(&addressForm); err != nil {
		api.HandleValidatorError(ctx, err)
		return
	}

	id := ctx.Param("id")
	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	_, err = global.AddressClient.UpdateAddress(context.Background(), &proto.AddressRequest{
		Id:           int32(i),
		Province:     addressForm.Province,
		City:         addressForm.City,
		District:     addressForm.District,
		Address:      addressForm.Address,
		SignerName:   addressForm.SignerName,
		SignerMobile: addressForm.SignerMobile,
	})
	if err != nil {
		zap.S().Errorw("更新地址失败")
		api.HandleGrpcErrorToHttp(err, ctx)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}
