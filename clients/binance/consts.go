package binance

// // SideType define side type of order
// type SideType string
//
// // OrderType define order type
// type OrderType string
//
// // TimeInForceType define time in force type of order
// type TimeInForceType string
//
// // NewOrderRespType define response JSON verbosity
// type NewOrderRespType string
//
// // OrderStatusType define order status type
// type OrderStatusType string
//
// // SymbolType define symbol type
// type SymbolType string
//
// // SymbolStatusType define symbol status type
// type SymbolStatusType string
//
// // SymbolFilterType define symbol filter type
// type SymbolFilterType string
//
// // UserDataEventType define spot user data event type
// type UserDataEventType string
//
// // MarginTransferType define margin transfer type
// type MarginTransferType int
//
// // MarginLoanStatusType define margin loan status type
// type MarginLoanStatusType string
//
// // MarginRepayStatusType define margin repay status type
// type MarginRepayStatusType string
//
// // FuturesTransferStatusType define futures transfer status type
// type FuturesTransferStatusType string
//
// // SideEffectType define side effect type for orders
// type SideEffectType string
//
// // FuturesTransferType define futures transfer type
// type FuturesTransferType int
//
// // TransactionType define transaction type
// type TransactionType string
//
// // LendingType define the type of lending (flexible saving, activity, ...)
// type LendingType string
//
// // StakingProduct define the staking product (locked staking, flexible defi staking, locked defi staking, ...)
// type StakingProduct string
//
// // StakingTransactionType define the staking transaction type (subscription, redemption, interest)
// type StakingTransactionType string
//
// // LiquidityOperationType define the type of adding/removing liquidity to a liquidity pool(COMBINATION, SINGLE)
// type LiquidityOperationType string
//
// // SwappingStatus define the status of swap when querying the swap history
// type SwappingStatus int
//
// // LiquidityRewardType define the type of reward we'd claim
// type LiquidityRewardType int
//
// // RewardClaimStatus define the status of claiming a reward
// type RewardClaimStatus int
//
// // RateLimitType define the rate limitation types
// // see https://github.com/binance/binance-spot-api-docs/blob/master/rest-api.md#enum-definitions
// type RateLimitType string
//
// // RateLimitInterval define the rate limitation intervals
// type RateLimitInterval string
//
// // Global enums
// const (
// 	SideTypeBuy  SideType = "BUY"
// 	SideTypeSell SideType = "SELL"
//
// 	OrderTypeLimit           OrderType = "LIMIT"
// 	OrderTypeMarket          OrderType = "MARKET"
// 	OrderTypeLimitMaker      OrderType = "LIMIT_MAKER"
// 	OrderTypeStopLoss        OrderType = "STOP_LOSS"
// 	OrderTypeStopLossLimit   OrderType = "STOP_LOSS_LIMIT"
// 	OrderTypeTakeProfit      OrderType = "TAKE_PROFIT"
// 	OrderTypeTakeProfitLimit OrderType = "TAKE_PROFIT_LIMIT"
//
// 	TimeInForceTypeGTC TimeInForceType = "GTC"
// 	TimeInForceTypeIOC TimeInForceType = "IOC"
// 	TimeInForceTypeFOK TimeInForceType = "FOK"
//
// 	NewOrderRespTypeACK    NewOrderRespType = "ACK"
// 	NewOrderRespTypeRESULT NewOrderRespType = "RESULT"
// 	NewOrderRespTypeFULL   NewOrderRespType = "FULL"
//
// 	OrderStatusTypeNew             OrderStatusType = "NEW"
// 	OrderStatusTypePartiallyFilled OrderStatusType = "PARTIALLY_FILLED"
// 	OrderStatusTypeFilled          OrderStatusType = "FILLED"
// 	OrderStatusTypeCanceled        OrderStatusType = "CANCELED"
// 	OrderStatusTypePendingCancel   OrderStatusType = "PENDING_CANCEL"
// 	OrderStatusTypeRejected        OrderStatusType = "REJECTED"
// 	OrderStatusTypeExpired         OrderStatusType = "EXPIRED"
//
// 	SymbolTypeSpot SymbolType = "SPOT"
//
// 	SymbolStatusTypePreTrading   SymbolStatusType = "PRE_TRADING"
// 	SymbolStatusTypeTrading      SymbolStatusType = "TRADING"
// 	SymbolStatusTypePostTrading  SymbolStatusType = "POST_TRADING"
// 	SymbolStatusTypeEndOfDay     SymbolStatusType = "END_OF_DAY"
// 	SymbolStatusTypeHalt         SymbolStatusType = "HALT"
// 	SymbolStatusTypeAuctionMatch SymbolStatusType = "AUCTION_MATCH"
// 	SymbolStatusTypeBreak        SymbolStatusType = "BREAK"
//
// 	SymbolFilterTypeLotSize          SymbolFilterType = "LOT_SIZE"
// 	SymbolFilterTypePriceFilter      SymbolFilterType = "PRICE_FILTER"
// 	SymbolFilterTypePercentPrice     SymbolFilterType = "PERCENT_PRICE"
// 	SymbolFilterTypeMinNotional      SymbolFilterType = "MIN_NOTIONAL"
// 	SymbolFilterTypeIcebergParts     SymbolFilterType = "ICEBERG_PARTS"
// 	SymbolFilterTypeMarketLotSize    SymbolFilterType = "MARKET_LOT_SIZE"
// 	SymbolFilterTypeMaxNumAlgoOrders SymbolFilterType = "MAX_NUM_ALGO_ORDERS"
//
// 	UserDataEventTypeOutboundAccountPosition UserDataEventType = "outboundAccountPosition"
// 	UserDataEventTypeBalanceUpdate           UserDataEventType = "balanceUpdate"
// 	UserDataEventTypeExecutionReport         UserDataEventType = "executionReport"
// 	UserDataEventTypeListStatus              UserDataEventType = "ListStatus"
//
// 	MarginTransferTypeToMargin MarginTransferType = 1
// 	MarginTransferTypeToMain   MarginTransferType = 2
//
// 	FuturesTransferTypeToFutures FuturesTransferType = 1
// 	FuturesTransferTypeToMain    FuturesTransferType = 2
//
// 	MarginLoanStatusTypePending   MarginLoanStatusType = "PENDING"
// 	MarginLoanStatusTypeConfirmed MarginLoanStatusType = "CONFIRMED"
// 	MarginLoanStatusTypeFailed    MarginLoanStatusType = "FAILED"
//
// 	MarginRepayStatusTypePending   MarginRepayStatusType = "PENDING"
// 	MarginRepayStatusTypeConfirmed MarginRepayStatusType = "CONFIRMED"
// 	MarginRepayStatusTypeFailed    MarginRepayStatusType = "FAILED"
//
// 	FuturesTransferStatusTypePending   FuturesTransferStatusType = "PENDING"
// 	FuturesTransferStatusTypeConfirmed FuturesTransferStatusType = "CONFIRMED"
// 	FuturesTransferStatusTypeFailed    FuturesTransferStatusType = "FAILED"
//
// 	SideEffectTypeNoSideEffect SideEffectType = "NO_SIDE_EFFECT"
// 	SideEffectTypeMarginBuy    SideEffectType = "MARGIN_BUY"
// 	SideEffectTypeAutoRepay    SideEffectType = "AUTO_REPAY"
//
// 	TransactionTypeDeposit  TransactionType = "0"
// 	TransactionTypeWithdraw TransactionType = "1"
// 	TransactionTypeBuy      TransactionType = "0"
// 	TransactionTypeSell     TransactionType = "1"
//
// 	LendingTypeFlexible LendingType = "DAILY"
// 	LendingTypeFixed    LendingType = "CUSTOMIZED_FIXED"
// 	LendingTypeActivity LendingType = "ACTIVITY"
//
// 	LiquidityOperationTypeCombination LiquidityOperationType = "COMBINATION"
// 	LiquidityOperationTypeSingle      LiquidityOperationType = "SINGLE"
//
// 	timestampKey  = "timestamp"
// 	signatureKey  = "signature"
// 	recvWindowKey = "recvWindow"
//
// 	StakingProductLockedStaking       = "STAKING"
// 	StakingProductFlexibleDeFiStaking = "F_DEFI"
// 	StakingProductLockedDeFiStaking   = "L_DEFI"
//
// 	StakingTransactionTypeSubscription = "SUBSCRIPTION"
// 	StakingTransactionTypeRedemption   = "REDEMPTION"
// 	StakingTransactionTypeInterest     = "INTEREST"
//
// 	SwappingStatusPending SwappingStatus = 0
// 	SwappingStatusDone    SwappingStatus = 1
// 	SwappingStatusFailed  SwappingStatus = 2
//
// 	RewardTypeTrading   LiquidityRewardType = 0
// 	RewardTypeLiquidity LiquidityRewardType = 1
//
// 	RewardClaimPending RewardClaimStatus = 0
// 	RewardClaimDone    RewardClaimStatus = 1
//
// 	RateLimitTypeRequestWeight RateLimitType = "REQUEST_WEIGHT"
// 	RateLimitTypeOrders        RateLimitType = "ORDERS"
// 	RateLimitTypeRawRequests   RateLimitType = "RAW_REQUESTS"
//
// 	RateLimitIntervalSecond RateLimitInterval = "SECOND"
// 	RateLimitIntervalMinute RateLimitInterval = "MINUTE"
// 	RateLimitIntervalDay    RateLimitInterval = "DAY"
// )
