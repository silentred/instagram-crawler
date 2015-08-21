package main


type Picture struct{
    Id string   `db:"orig_id"`
    Url string  `db:"pic_url"`
    Status int   `db:"status"`
    CreatedTime int `db:"created_time"`
}

type User struct{
    Id string   `db:"orig_id"`
    Name string `db:"name"`
    AccessToken string  `db:"access_token"`
    LastAuthTime int    `db:"last_auth_time"`
    Valid int   `db:"valid"`
}
