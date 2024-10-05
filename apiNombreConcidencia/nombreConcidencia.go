package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)
  var db *sql.DB

  //Creamos nuestra estructura
  type Emprendedor struct{
	Id int `json:"Id"`
	Correo  string  `json:"Correo"`
	Negocio DescripcionNegocio `json:"descripcionNegocio"`
  } 
  type DescripcionNegocio struct{
	Id_negocio int `json:"Id_negocio"`
	NombreNegocio string `json:"NombreNegocio"`
	Direccion string `json:"Direccion"`
	Descripcion string `json:"Descripcion"`
	Verificacion string `json:"Verificacion"`
  }
  //Creamos nuestra funcion para conectarnos a la base de datos y tener una mejor autenticacion
  func initDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{
		log.Fatal("Error en la conexion en la base de  datos intente lo de nuevo")
	}

	//Verificamos si la conexion esta activa mediante  un pin
	err=db.Ping()
	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("Se ha conectado correctamente a la base de datos BizzBuz")
  }
func main() {
	 initDB()
	 router:=mux.NewRouter()
	 router.HandleFunc("/busqueda",getBusquedaCoincidencia).Methods("GET")
	 headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	 originOk:=handlers.AllowedOrigins([]string{"*"})
	 methodsOk:=handlers.AllowedMethods([]string{"GET,HEAD,POST,PUT,OPTIONS"})
	 log.Fatal(http.ListenAndServe(":5700",handlers.CORS(headerOk,originOk,methodsOk)(router)))

}

func getBusquedaCoincidencia(w http.ResponseWriter,r *http.Request){
      
	nombre:= r.URL.Query().Get("nombre")
	if nombre==""{
		log.Fatal(w,"El parametro 'nombre' es requerido para la busqueda",http.StatusBadRequest)
		return
	}

	//Realizamos nuestra consulta
	query:= `
       SELECT d.NombreNegocio,d.Direccion,d.Verificacion,d.Id_negocio,d.Descripcion, e.Id,e.Correo
	   FROM Emprendedor e INNER JOIN descripcionNegocio d ON e.id_negocio=d.Id_negocio
	   WHERE e.Nombre LIKE @p1 OR d.NombreNegocio LIKE @p2

	`
    rows,err:=db.Query(query,"%"+strings.TrimSpace(nombre)+"%","%"+strings.TrimSpace(nombre)+"%")
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	var emprendedores []Emprendedor
	for rows.Next(){
		var emprendedor Emprendedor
		var descripcion DescripcionNegocio
		err=rows.Scan(&descripcion.NombreNegocio,&descripcion.Direccion,&descripcion.Verificacion,&descripcion.Id_negocio,&descripcion.Descripcion,&emprendedor.Id,&emprendedor.Correo)
		if err !=nil{
			log.Fatal(w,err.Error(),http.StatusInternalServerError)
			return
		}
            emprendedor.Negocio=descripcion
	         emprendedores = append(emprendedores, emprendedor)
		}

		w.Header().Set("Content-Type","application/json")
		json.NewEncoder(w).Encode(emprendedores)
}  