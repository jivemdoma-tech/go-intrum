package gointrum

type WebhookEvent struct {
	SubjectType   subjectType `json:"subject_type"`
	SubjectTypeID int64       `json:"subject_type_id,string"`
	Event         event       `json:"event"`
	ObjectType    objectType  `json:"object_typ"`
	ObjectSubType int64       `json:"object_sub_type,string"`
	ObjectSubID   int64       `json:"object_sub_id,string"`
	Snapshot      Snapshot    `json:"snapshot"`
}
type Snapshot struct {
	Merge any `json:"merge"`
}
