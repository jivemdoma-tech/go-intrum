package gointrum

type subjectType string

const (
	SubjectTypeSystem   subjectType = "system"
	SubjectTypeBusiness subjectType = "business"
	SubjectTypeEmployee subjectType = "employee"
)

type event string

const (
	EventLogin     event = "login"
	EventView      event = "view"
	EventCreate    event = "create"
	EventEdit      event = "edit"
	EventDelete    event = "delete"
	EventExport    event = "export"
	EventImport    event = "import"
	EventStatus    event = "status"
	EventStage     event = "stage"
	EventComment   event = "comment"
	EventManager   event = "manager"
	EventAnswer    event = "answer"
	EventMessenger event = "messenger"
	EventQueue     event = "queue"
	EventOther     event = "other"
)

type objectType string

const (
	ObjectTypeCustomer    objectType = "customer"
	ObjectTypeRequest     objectType = "request"
	ObjectTypeStock       objectType = "stock"
	ObjectTypeSale        objectType = "sale"
	ObjectTypeTask        objectType = "task"
	ObjectTypeMessenger   objectType = "messenger"
	ObjectTypeRemind      objectType = "remind"
	ObjectTypeEmail       objectType = "email"
	ObjectTypeEmailsystem objectType = "emailsystem"
	ObjectTypeCall        objectType = "call"
	ObjectTypeSms         objectType = "sms"
	ObjectTypeDelivery    objectType = "delivery"
	ObjectTypeComment     objectType = "comment"
	ObjectTypeBlank       objectType = "blank"
	ObjectTypeApp         objectType = "app"
)
