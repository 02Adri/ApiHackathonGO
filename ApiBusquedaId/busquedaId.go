package main 

import(
	"database/sql"
	"encoding/json"
	"net/http"
	"log"
	"fmt"
_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)


//creamos nuestra estructuras
type Emprendedor struct{
   Id int `json:"Id"`
   Nombre string `json:"Nombre"`
   Correo  string `json:"Correo"`
   Password string `json:"Password"`
   Telefono string `json:"Telefono"`
   Id_negocio string `json:"id_negocio"`
    Description DescripcionNegocio `json:"descripcionNegocio"`
}

type DescripcionNegocio struct{

	Id_negocio int `json:"Id_negocio"`
	NombreNegocio string `json:"NombreNegocio"`
	Direccion string `json:"Direccion"`
	Descripcion string `json:"Descripcion"`
	Verificacion string `json:"Verificacion"`
}
  //primero inicializamos nuestras variables
  var db *sql.DB
  
  //creamos la funcion para conectar a la base de datos
  func initDB(){

	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	 db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{

		log.Fatal("No se pudo conectar a la base de datos BizzBuzz intente lo de nuevo")
	}
     //Verificamos que la conexion este activa
	err=db.Ping()
    if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("Se ha conectado correctamente a la base de datos BizzBuzz")
  }
func main(){
   
//mandamos a llamar nuestra funcion para conectar nuestra base de datos
initDB()
  router:=mux.NewRouter()
  router.HandleFunc("/busquedaId/{Id}",getBusquedaId).Methods("GET")
  headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
  originOk:=handlers.AllowedOrigins([]string{"*"})
  methodsOk:=handlers.AllowedMethods([]string{"GET,HEAD,POST,PUT,OPTIONS"})
  log.Fatal(http.ListenAndServe(":5800",handlers.CORS(headerOk,originOk,methodsOk)(router)))

}

func getBusquedaId(w http.ResponseWriter,r *http.Request){
   
	vars:=mux.Vars(r)
	Id:=vars["Id"]

	 rows,err:=db.Query(`SELECT e.Id,e.Nombre,e.Correo,e.Password,e.Telefono,e.id_negocio,d.Id_Negocio,d.NombreNegocio,d.Direccion,
	 d.Descripcion,d.Verificacion FROM Emprendedor e INNER JOIN descripcionNegocio d ON e.id_negocio=d.Id_negocio 
	 WHERE e.Id=@p1`,Id)

	if err !=nil{
       log.Fatal(w,err.Error(),http.StatusInternalServerError)
		 return
	 }

	 var Emprendedores []Emprendedor
	 for rows.Next(){
		var Emprendedor Emprendedor
		var Descripcion DescripcionNegocio
         err=rows.Scan(&Emprendedor.Id,&Emprendedor.Nombre,&Emprendedor.Correo,&Emprendedor.Password,&Emprendedor.Telefono,&Emprendedor.Id_negocio,&Descripcion.Id_negocio,&Descripcion.NombreNegocio,&Descripcion.Direccion,&Descripcion.Descripcion,&Descripcion.Verificacion)
		 if err !=nil{
			log.Fatal(w,err.Error(),http.StatusInternalServerError)
			return
		 }

		 Emprendedor.Description=Descripcion

		 Emprendedores=append(Emprendedores, Emprendedor)
	 }
	 w.Header().Set("Content-Type","aplication/json")
	 json.NewEncoder(w).Encode(Emprendedores)
}