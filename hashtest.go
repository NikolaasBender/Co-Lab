package main

import (
  "golang.org/x/crypto/bcrypt"
  "fmt"
)


func main() {
  pss1 := "1234"
  // pss2 :=
  // pss3 :=
  // pss4 :=
  // pss5 :=
  // pss6 :=
  // pss7 :=
  // pss8 :=
  // pss9 :=
  // pss10 :=
  // pss11 :=
  // pss12 :=
  // pss13 :=
  // pss14 :=
  // pss15 :=
  // pss16 :=


  pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss1),bcrypt.MinCost)
  out := string(pss_hash)
  fmt.Println(out)

  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss2),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss3),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss4),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss5),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss6),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss7),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss8),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss9),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss10),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss11),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss12),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss13),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss14),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss15),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
  //
  // pss_hash, _ := bcrypt.GenerateFromPassword([]byte(pss16),bcrypt.MinCost)
  // out := string(pss_hash)
  // fmt.Println(out)
}
