package main
import (
	"database/sql"
	"encoding/json"
	"net/http"
	"log"
	"fmt"//enviar mensaje
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Realizamos una asignacion a la base de datos
var db  *sql.DB
//Creamos la estructura de nuestra instancia
type Cliente struct{
 Id int `json:"Id"`
 Nombre string `json:"Nombre"`
 Correo string `json:"Correo"`
 Password string `json:"Password"`
 Telefono string `json:"Telefono"`

}

//Realizamos una funcion para iniciar nuestra base de datos
func initDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err!=nil{
		log.Fatal("Error al conectarse a la base de datos intentelo lo de nuevo")
	}

	err=db.Ping()//realizamos un ping para visualizar si hay una conexion activa
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Conexion Exitosa a la base de datos...")
}
 
func main() {
	initDB()
	router:=mux.NewRouter()
	router.HandleFunc("/Cliente/{Id}",getCliente).Methods("GET")
	router.HandleFunc("/Cliente",postCliente).Methods("POST")
	//Realizamos los permisos de consumo de datos
	headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET,HEAD,POST,PUT,OPTIONS"})
	log.Fatal(http.ListenAndServe(":3600",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}
func getCliente(w http.ResponseWriter, r *http.Request){

	vars:=mux.Vars(r)
	Id:= vars["Id"]
   rows,err:=db.Query("SELECT Id,Nombre,Correo,Password,Telefono FROM Cliente WHERE Id=@p1",Id)
   if err !=nil{
	log.Fatal(w,err.Error(),http.StatusInternalServerError)
	return
   }
   var clientes []Cliente
   for rows.Next(){
     var cliente Cliente
	 err=rows.Scan(&cliente.Id,&cliente.Nombre,&cliente.Correo,&cliente.Password,&cliente.Telefono)
	 if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	   }
	 clientes=append(clientes, cliente)
   }
   w.Header().Set("Content-Type","encoding/json")
   json.NewEncoder(w).Encode(clientes)
}
func postCliente(w http.ResponseWriter, r *http.Request){
	var cliente Cliente
	err:=json.NewDecoder(r.Body).Decode(&cliente)
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
		return
	}

	if cliente.Nombre=="" ||cliente.Correo==""||cliente.Password==""||cliente.Telefono==""{
		log.Fatal(w,"No se permiten campos nulos, ingreselos correctamente",http.StatusBadRequest)
		return
	}
	_,err=db.Query("Insert INTO Cliente(Nombre,Correo,Password,Telefono) VALUES(@Nombre,@Correo,@Password,@Telefono)",
       sql.Named("Nombre",cliente.Nombre),
	   sql.Named("Correo",cliente.Correo),
	   sql.Named("Password",cliente.Password),
	   sql.Named("Telefono",cliente.Telefono),)
     if err !=nil {
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	 }
	 w.WriteHeader(http.StatusCreated)
	 fmt.Fprintf(w,"Cliente registrado Correctamente")
}