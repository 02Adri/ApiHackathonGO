package main

import (
	"database/sql"
	/*"encoding/base64"
	"encoding/json"*/
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//realizamos nuestra direccion de puntero de nuestra base de datos
  var db *sql.DB

//realizamos una constante la cual te permita guardar la carpeta de las imagenes
const (
	uploadsPath="./uploads/"
)
//funcion para inicial la base de datos
func initDB(){
   var err error
   connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
   db,err=sql.Open("sqlserver",connectionString)
    if err !=nil{
		log.Fatal("Error al conectarse a la base de datos con la imagen")
	}

	err=db.Ping()//crear un ping para ver si esta activa la base de datos
	if err!=nil{

		log.Fatal(err)
	}

	fmt.Println("Conectado correctamente a la base de datos de la imagen")

}

func main() {
	  initDB()
	   router:=mux.NewRouter()
	   router.HandleFunc("/imagenEmprendedor/{Id}",getImagenEmprendedor).Methods("GET")
	   router.HandleFunc("/imagenEmprendedor",postImagenEmprendedor).Methods("POST")
	   headerOk:=handlers.AllowedHeaders([]string{"X-Request-With","Content-Type","Authorization"})
	   originOk:=handlers.AllowedOrigins([]string{"*"})
	   methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
	   log.Fatal(http.ListenAndServe(":5100",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

func getImagenEmprendedor(w http.ResponseWriter, r *http.Request){
  
	//Obtenemos el valor del id para hacer referencia a la sustraccion de la imagen
	 vars:=mux.Vars(r)
	 Id:= vars["Id"]

    var fileName string

	err:=db.QueryRow("SELECT nombreImagen FROM imagenEmprendedor WHERE Id=@p1",Id).Scan(&fileName)

    if err!=nil{
		if err==sql.ErrNoRows{
			http.NotFound(w,r)//Si no encuentra en el cuerpo de l base de datos
			return
		}else{
			log.Fatal(w,err.Error(),http.StatusInternalServerError)
			return
		}
	}
     filePath:=filepath.Join(uploadsPath,fileName)

	 if _,err:=os.Stat(filePath); os.IsNotExist(err){
		http.NotFound(w,r)
		return
	 }
	 http.ServeFile(w,r,filePath)
}

func postImagenEmprendedor(w http.ResponseWriter, r *http.Request){

	 err:=r.ParseMultipartForm(10<<20)//debe ser menor a 10mb file
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	file,fileHeader,err:=r.FormFile("image")
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	defer file.Close()

	Id:=r.FormValue("Id")
	formExt:=filepath.Ext(fileHeader.Filename)
	fileName:=fmt.Sprintf("%s%s",Id,formExt)
	filePath:=filepath.Join("uploads",fileName)

	//creamos el path
	out,err:=os.Create(filePath)
	if err !=nil{
		log.Fatal(w,"error al enviar el formato de imagen",http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_,err=io.Copy(out,file)
	if err !=nil{
		log.Fatal(w,"Imagen invalida no es el archivo correspondiente",http.StatusInternalServerError)
		return
	}

	_,err=db.Exec(`
	MERGE INTO imagenEmprendedor AS target
	  USING(SELECT @p1 AS Id, @p2 AS nombreImagen)
	  AS source ON (target.Id=source.Id) 
	  WHEN MATCHED THEN UPDATE SET nombreImagen=source.nombreImagen
	  WHEN NOT MATCHED THEN INSERT(Id,nombreImagen) VALUES(source.Id,source.nombreImagen);
	`,Id,fileName)

	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("La imagen se envio Correctamente"))
}