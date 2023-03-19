package constants

var EventTypeNumbers = map[string]int{
    "customer_created": 1,
    "customer_action_performed": 2,
    "customer_billed": 3,
    "customer_deactivated": 4,
}

var StatusNumbers = map[string]int{
    "success": 100,
    "error": 200,
}