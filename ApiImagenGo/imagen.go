package main

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Image struct{
	Id_image int `json:"id_image"`
	Nombre   string `json:"Nombre"`
	Data string `json:"Data"`
}

//Realizamos la conexion a la base de datos 
func initDB(){
    var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=ApiLoginSismo"
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{
		log.Fatal("Error en la conexion a la base de datos intentelo  de nuevo")

	}
	err=db.Ping()//realizamos un ping para saber si esta activa la conexion
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Println("Conectado Correctamente a la base de datos...")
}

func main() {
	//Llamamos  la funcion para la conexion a la base de datos
	initDB()
	//inicializamos nuestra ruta
	router:=mux.NewRouter()
	//Llamamos nuestros metodos
	router.HandleFunc("/imagen/{id}",getImagen).Methods("GET")
	router.HandleFunc("/imagen",postImagen).Methods("POST")
	headerOk:=handlers.AllowedHeaders([]string{"X-Request-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
	log.Fatal(http.ListenAndServe(":2700",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

//Realizamos la funcion para enviar la imagen
func getImagen(w http.ResponseWriter, r *http.Request){
   
    params:=mux.Vars(r)
	id:=params["id_image"]
	var image Image
	var imageData []byte
	
	err:= db.QueryRow("SELECT id_image,Nombre,Data FROM Imagen WHERE id_image=@p1",id).Scan(&image.Id_image,&image.Nombre,&imageData)
     if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
	    return
	}
    image.Data=base64.RawURLEncoding.EncodeToString(imageData)
   
      json.NewEncoder(w).Encode(image)
}

 func postImagen(w http.ResponseWriter, r *http.Request){
	
	err:=r.ParseMultipartForm(10<<20)//debe ser menor a 10 mb file
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
	      return
		} 

	file,handler,err:=r.FormFile("file")
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fileBytes,err:=io.ReadAll(file)
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
        
	//name:=r.FormValue("Nombre")
	filePath:="./uploads/"+handler.Filename
	err=os.WriteFile(filePath,fileBytes,0644)
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
	_,err=db.Query("INSERT INTO Imagen(Nombre,Data)VALUES(@p1,@p2)",handler.Filename,fileBytes)
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w,"Imagen Enviada correctamenten:%s\n",handler.Filename)

}