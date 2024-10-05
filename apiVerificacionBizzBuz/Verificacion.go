package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
_ "github.com/denisenkom/go-mssqldb"
"github.com/gorilla/handlers"
 "github.com/gorilla/mux"
)

 var db *sql.DB

 type   User struct{
	Id int `json:"ID"`
	Nombre string `json:"Nombre"`
   Correo string `json:"Correo"`
   Password string `json:"Password"`
   Telefono string `json:"Telefono"`
   Table string `json:"table"`
   Id_negocio int `json:"id_negocio"`
 }

 func initDB(){
  var err error
   connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
   db,err=sql.Open("sqlserver",connectionString)
   if err !=nil{
       log.Fatal("Error en la conexion a la base de datos intente lo de nuevo")
   }
   err=db.Ping()//verificamos que la conexion de la base de datos este activa

   if err !=nil{
	  log.Fatal(err)
   }
   fmt.Println("Conexion exitosa a la base de datos BizzBuz")
 }
func main() {
	initDB()
	router:=mux.NewRouter()
	router.HandleFunc("/verificacionUser",postVerificacionUser).Methods("POST")
	headerOk:=handlers.AllowedHeaders([]string{"X-Request-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET,HEAD,POST,PUT,OPTIONS"})
	log.Fatal(http.ListenAndServe(":5300",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}
  
 
func postVerificacionUser(w http.ResponseWriter, r *http.Request){
       
	var user User
	_ =json.NewDecoder(r.Body).Decode(&user)
	
	
	if user.Table=="cliente"{
		
		query:="SELECT Id,Nombre,Correo,Password,Telefono FROM Cliente WHERE Correo=@p1 AND Password=@p2"
	err:=db.QueryRow(query,user.Correo,user.Password).Scan(&user.Id,&user.Nombre,&user.Correo,&user.Password,&user.Telefono)
	if err!=nil{
	
		fmt.Println("El usuario cliente no se encuentra registrado en la aplicacion BizzBuz")
		return
	}

	}else if user.Table=="emprendedor"{
	
		query:="SELECT Id,Nombre,Correo,Password,Telefono,id_negocio FROM Emprendedor WHERE Correo=@p1 AND Password=@p2"
	err:=db.QueryRow(query,user.Correo,user.Password).Scan(&user.Id,&user.Nombre,&user.Correo,&user.Password,&user.Telefono,&user.Id_negocio)
	if err!=nil{
	
		fmt.Println("El usuario  emprendedor no se encuentra registrado en la aplicacion BizzBuz")
		return
	}
	}else{
		log.Fatal(w,"Parametros de la tabla no validos",http.StatusBadRequest)
		return
	}

	/*var id int
	query:="SELECT IdFROM "+tableName+" WHERE Correo=@p1 AND Password=@p2"
	err:=db.QueryRow(query,user.Correo,user.Password).Scan(&id)
	if err!=nil{
	
		fmt.Println("El usuario no se encuentra registrado en la aplicacion BizzBuz")
		return
	}
	*/
    
		response:=map[string]interface{}{
			"message":"El usuario se encuentra registrado en la aplicacion BizzBuz",
			"user":user,
		}
	 w.Header().Set("Content-Type","application/json")
	 json.NewEncoder(w).Encode(response)
	}
