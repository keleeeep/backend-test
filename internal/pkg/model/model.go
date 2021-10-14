/*
 * @Author: Adrian Faisal
 * @Date: 14/10/21 13.41
 */

package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	Role      string    `json:"role"`
	Timestamp time.Time `json:"timestamp"`
}

type Data struct {
	Uuid         string    `json:"uuid"`
	Komoditas    string    `json:"komoditas"`
	AreaProvinsi string    `json:"area_provinsi"`
	AreaKota     string    `json:"area_kota"`
	Size         string    `json:"size"`
	Price        string    `json:"price"`
	Usd          float64   `json:"usd"`
	TglParsed    time.Time `json:"tgl_parsed"`
	Timestamp    string    `json:"timestamp"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (p *User) TableName() string {
	return "users"
}
