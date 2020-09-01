package models

type CommonReturn struct {
	Code int
	Msg  string
}

func (c *CommonReturn) SetData(code int, msg string) {
	c.Code, c.Msg = code, msg
}

type TradeAccountSingReq struct {
	AccountID uint64 `json:"accountId"`
}
type TradeAccountArrReq struct {
	AccountIDs []uint64 `json:"accountId"`
}

type TradAccountData struct {
	ID           int32
	AccountID    uint64 `protobuf:"varint,1,opt,name=accountID,proto3" json:"accountID,omitempty"`
	NickName     string `protobuf:"bytes,2,opt,name=nickName,proto3" json:"nickName,omitempty"`
	Userid       string `protobuf:"bytes,3,opt,name=userid,proto3" json:"userid,omitempty"`
	Leverage     uint64 `protobuf:"varint,4,opt,name=leverage,proto3" json:"leverage,omitempty"`
	Currency     string `protobuf:"bytes,5,opt,name=currency,proto3" json:"currency,omitempty"`
	ReadPassword string `protobuf:"bytes,6,opt,name=readPassword,proto3" json:"readPassword,omitempty"`
	Password     string `protobuf:"bytes,7,opt,name=password,proto3" json:"password,omitempty"`
	LastLogin    uint64 `protobuf:"varint,8,opt,name=lastLogin,proto3" json:"lastLogin,omitempty"`
	IsActive     bool   `protobuf:"varint,9,opt,name=isActive,proto3" json:"isActive,omitempty"`
	IsDelete     bool   `protobuf:"varint,10,opt,name=isDelete,proto3" json:"isDelete,omitempty"`
	DateJoined   uint64 `protobuf:"varint,11,opt,name=dateJoined,proto3" json:"dateJoined,omitempty"`
	Email        string `protobuf:"bytes,12,opt,name=email,proto3" json:"email,omitempty"`
	CountryCode  int32  `protobuf:"varint,13,opt,name=countryCode,proto3" json:"countryCode,omitempty"`
	Phone        string `protobuf:"bytes,14,opt,name=phone,proto3" json:"phone,omitempty"`
}

type TradAccountUpdateReq struct {
	Data TradAccountData
}
