package types

type ModelResult struct {
	ID           int
	ModelID      int
	JobID        string
	SourceID     string
	JobSourceID  int
	CreatedAt    string
	Registration string
	Images       string
	Gearbox      string
	Mileage      int
	Price        int
	Fuel         string
	Description  string
	EngineSize   string
	Spec         string
	Owners       int
	URL          string
	Seller       string
	Year         int
	Power        string
	Cat          string
	Colour       string
	Doors        int
	DriveType    string
	BodyType     string
	SaleType     string
}

// func CreateModelResult(v *Vehicle) *ModelResult {
// 	return &ModelResult{
// 		SourceID:     v.SourceID,
// 		Registration: v.Registration,
// 		Images:       *v.Image.Main.Url,
// 		Gearbox:      v.Gearbox,
// 		Mileage:      v.Mileage,
// 		Price:        *v.Price,
// 		Fuel:         v.Fuel.Description,
// 		EngineSize:   v.EngineSize,
// 		Spec:         v.Spec,
// 		Year:         v.RegYear,
// 	}
// }

/*

create table main.model_result
(
    id            integer
        primary key autoincrement,
    created_at    TIMESTAMP default CURRENT_TIMESTAMP not null,
    registration  text,
    images        text,
    gearbox       text,
    mileage       bigint,
    price         bigint,
    fuel          text,
    description   text,
    engine_size   text,
    spec          text,
    owners        bigint,
    url           text,
    seller        text,
    year          bigint,
    model_id      integer
        references main.model,
    job_id        text,
    source_id     text,
    power         text,
    job_source_id bigint
        references main.job_source,
    cat           text,
    colour        text,
    doors         bigint,
    driveType     text,
    bodyType      text,
    sale_type     text
);

create index main.idx_model_results_job_id
    on main.model_result (job_id);

*/
