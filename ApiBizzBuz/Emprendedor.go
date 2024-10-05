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
var db *sql.DB//lo que permite es hacer una asignacion  estatica

//crear una estructura
type Emprendedor struct{
	Id  int  `json:"Id"`
	Nombre string `json:"Nombre"`
	Correo string `json:"Correo"`
	Password string `json:"Password"`
	Telefono  string `json:"Telefono"`
	Id_negocio int `json:"id_negocio"`
}

type descripcionNegocio struct{
	Id_Negocio int `json:"id_negocio"`
	NombreNegocio string `json:"NombreNegocio"`
	 Direccion string `json:"Direccion"`
	 Descripcion string `json:"Descripcion"`
	 Verificacion string `json:"Verificacion"` 
}

type EmprendedorDescripcionNegocio struct{
	EmprendedorC Emprendedor `json:"Emprendedor"`
	DescripcionNegocioC descripcionNegocio `json:"descripcionNegocio"`
}

func initDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{
		log.Fatal("Error en la conexion a la base de datos BizzBuz intentelo de nuevo")
	}
	err=db.Ping()
	if err !=nil{
		log.Fatal(err)
	}

	fmt.Println("Se ha conectado correctamente a la base de datos BizzBuz")
}

func  main()  {
	 initDB()//mando a llamar la funcion de la base de datos
   	router:= mux.NewRouter()
	router.HandleFunc("/Emprendedor/{id_negocio}",getEmprendedor).Methods("GET")
	router.HandleFunc("/Emprendedor",postEmprendedor).Methods("POST")
    headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	originOK:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
     log.Fatal(http.ListenAndServe(":3500",handlers.CORS(headerOk,originOK,methodsOk)(router)))

}
func getEmprendedor(w http.ResponseWriter,r *http.Request){
	 vars:=mux.Vars(r)
	 id_negocio:=vars["id_negocio"]

	 rows,err:=db.Query(`SELECT e.*, d.* FROM Emprendedor e JOIN descripcionNegocio d ON e.id_negocio=d.id_negocio
	 WHERE e.id_negocio=@p1`,id_negocio)
	 if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	 }

	 var emprenDescripNegocio []EmprendedorDescripcionNegocio
	 for rows.Next(){
		var emprendedor Emprendedor
		var descripcion descripcionNegocio
		err=rows.Scan(&emprendedor.Id,&emprendedor.Nombre,&emprendedor.Correo,&emprendedor.Password,&emprendedor.Telefono,&emprendedor.Id_negocio,&descripcion.Id_Negocio,&descripcion.NombreNegocio,&descripcion.Direccion,&descripcion.Descripcion,&descripcion.Verificacion)
		if err !=nil{
			log.Fatal(w,err.Error(),http.StatusInternalServerError)
			return
		 }
		 emprenDescripNegocio=append(emprenDescripNegocio, EmprendedorDescripcionNegocio{
			EmprendedorC:emprendedor,
			DescripcionNegocioC: descripcion,
		 })
		 w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(emprenDescripNegocio)
	 }
}
func postEmprendedor(w http.ResponseWriter, r *http.Request){

      var eDescripcion EmprendedorDescripcionNegocio
	err:=json.NewDecoder(r.Body).Decode(&eDescripcion)
	//err=json.NewDecoder(r.Body).Decode(&descripcionnegocio)
	 if err !=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
		return
	 }
	 //Realizamos la consulta del procedimiento almacenado
	_,err=db.Exec("EXEC emprendedorNegocio @Nombre=@Nombre, @Correo=@Correo, @Password=@Password,@Telefono=@Telefono,@NombreNegocio=@NombreNegocio,@Direccion=@Direccion,@Descripcion=@Descripcion,@Verificacion=@Verificacion",
	   sql.Named("Nombre",eDescripcion.EmprendedorC.Nombre),
	   sql.Named("Correo",eDescripcion.EmprendedorC.Correo),
	   sql.Named("Password",eDescripcion.EmprendedorC.Password),
	   sql.Named("Telefono",eDescripcion.EmprendedorC.Telefono),
	   sql.Named("NombreNegocio",eDescripcion.DescripcionNegocioC.NombreNegocio),
	   sql.Named("Direccion",eDescripcion.DescripcionNegocioC.Direccion),
	   sql.Named("Descripcion",eDescripcion.DescripcionNegocioC.Descripcion),
	   sql.Named("Verificacion",eDescripcion.DescripcionNegocioC.Verificacion),
       )
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	
}