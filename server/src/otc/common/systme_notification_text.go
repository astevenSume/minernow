package common

const (
	BuerPayedTitle = "买家已付款!"
	BuerPayedBody  = "您订单号为『%d』，数量为『%.4f』个，总价为『%.2f』CNY的订单对方已完成支付，请及时查账，确认到账后请及时放币，若未收到款项，可申请客诉处理。" //收款订单-买家已付款

	BuerOrderCancelTitle = "买家取消订单!"
	BuerOrderCancelBody  = "您的「收款」订单『%d』已被对方取消。"

	BuerReceiptedTitle = "买家确认收款!"
	BuerReceiptedBody  = "『%d』已确认收到您的订单『%d』付款，系统会自动将『%.4f』EUSD发放到您的账户，请注意查收。"

	SysterOrderCancelTitle = "系统取消订单!"
	SysterOrderCancelBody  = "您的订单『%d』因为超时已被系统取消。"

	ExchangerPayedTitle = "卖家已付款!"
	ExchangerPayedBody  = "您订单号为『%d』，数量为『%.4f』个，总价为『%.2f』CNY的订单对方已完成支付，请及时查账，确认到账后请及时放币，若未收到款项，可申请客诉处理。"

	ExchangerSendedTitle = "卖家已放币!"
	ExchangerSendedBody  = "『%d』已确认收到您的订单『%d』付款，系统会自动将您所购买的『%.4f』EUSD发放到您的账户，请注意查收。"

	ExchangerNewOrder1Title = "收款订单-新订单!"
	ExchangerNewOrder1Body  = "您有一笔新的收款订单待处理，订单号为『%d』，数量为『%.4f』个，总价为『%.2f』CNY，请等待对方付款后及时查账确认。"

	ExchangerNewOrder2Title = "支付订单-新订单!"
	ExchangerNewOrder2Body  = "您有一笔新的支付订单待处理，订单号为『%d』，数量为『%.4f』个，总价为『%.2f』CNY，请等待对方转币后及时查账确认。"
)
const (
	NewMessage = "订单聊天-新消息"
)
