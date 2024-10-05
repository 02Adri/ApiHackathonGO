package main

import(
	"database/sql"
	"encoding/json"
	"net/http"
	"fmt"
	"log"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// incializamos nuestro puntero
 var db *sql.DB
//Creamos nuestra estructura 
type Producto struct{
	Id int `json:"Id"`
	Nombre string `json:"Nombre"`
	Categoria string `json:"Categoria"`
	Descripcion string `json:"Descripcion"`
	Precio float64 `json:"Precio"`
}

//Creamos la funcion para poder conectarnos a la base de datos
func InitDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil {
		log.Fatal("Error al conectarse a la base de datos intente lo nuevamente")
	}

	err=db.Ping()// Comprobamos que este activa la conexion a la base de datos
	if err !=nil{
		log.Fatal(err)
	}
    fmt.Println("Se ha conectado correctamente a la base de datos BizzBuzz")
}
func main() {
	//llamamos nuestra funcion para conectarse a l base de datos
     InitDB()
	 router:=mux.NewRouter();
	 router.HandleFunc("/Producto",getProducto).Methods("GET")
	 router.HandleFunc("/Producto",postProducto).Methods("POST")
	 headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	 originOk:=handlers.AllowedOrigins([]string{"*"})
	 methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
	 log.Fatal(http.ListenAndServe(":5300",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

func getProducto(w http.ResponseWriter, r *http.Request){
    var err error
	rows,err:=db.Query("SELECT Id,Nombre,Categoria,Descripcion,Precio FROM Producto")

	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
	defer rows.Close()//Cerramos la conexion para la consulta

	var productos []Producto
   for rows.Next(){
	var producto Producto
	err=rows.Scan(&producto.Id,&producto.Nombre,&producto.Categoria,&producto.Descripcion,&producto.Precio)
	if err!=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
     productos=append(productos, producto)
   }
   w.Header().Set("Content-Type", "encoding/json")
   json.NewEncoder(w).Encode(productos)
}

func postProducto(w http.ResponseWriter,r *http.Request){
  var productos Producto
  err:= json.NewDecoder(r.Body).Decode(&productos)

  if  err!=nil{
	 log.Fatal(w,err.Error(),http.StatusBadRequest)
	 return
  }

  if productos.Nombre==""|| productos.Categoria==""|| productos.Descripcion==""||productos.Precio==0{
	log.Fatal(w,"Ingrese todos los campos porque son obligatorios",http.StatusBadRequest)
	return
  }

  _,err=db.Query("INSERT Producto(Nombre,Categoria,Descripcion,Precio) VALUES(@Nombre,@Categoria,@Descripcion,@Precio)",
    sql.Named("Nombre",productos.Nombre),
	sql.Named("Categoria",productos.Categoria),
	sql.Named("Descripcion",productos.Descripcion),
	sql.Named("Precio",productos.Precio),
)

  if err!=nil{

	log.Fatal(w,err.Error(),http.StatusInternalServerError)
	return
  }
  w.WriteHeader(http.StatusCreated)
  w.Write([]byte("Producto registrado correctamente en la  empresa del emprendedor"))
}