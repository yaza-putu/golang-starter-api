package validation

type Category struct {
	Name string `json:"name" form:"name" validate:"required,unique=categories:name"`
}

type CategoryWithIgnore struct {
	ID   string `json:"id" form:"id"`
	Name string `json:"name" form:"name" validate:"required,unique=categories:name:ID"`
}
