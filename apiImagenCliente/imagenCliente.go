package main

import (
	"database/sql"
	
    "fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
      "io"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

//Realizamos nuestro puntero de sql
   var db *sql.DB
    const(
		uploadPath="./uploads/"
	)
   

//funcion para inicializar nuestra base de datos
func initDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err !=nil{
		log.Fatal("Error en la conexion en la base de datos con la imagen cliente")
	}
	err=db.Ping()//realizamos un ping de conexion con la base de datos
	if err !=nil{

		log.Fatal(err)
	}

	fmt.Println("Se ha conectado correctamente la base de datos...")
}

func main() {
	initDB()
	router:=mux.NewRouter()
	router.HandleFunc("/imagenCliente/{Id}",getImagenCliente).Methods("GET")
	router.HandleFunc("/imagenCliente",postImagenCliente).Methods("POST")
	headerOk:=handlers.AllowedHeaders([]string{"X-Request-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET,HEAD,POST,PUT,OPTIONS"})
	log.Fatal(http.ListenAndServe(":5200",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

func getImagenCliente(w http.ResponseWriter, r *http.Request){
    
  vars:=mux.Vars(r)
  Id:=vars["Id"]
  //Obtenemos el archivo  de la base de datos
   var fileName string
   err:=db.QueryRow("SELECT nombreImagen FROM  imagenCliente WHERE Id=@p1 ",Id).Scan(&fileName)

    if err != nil{

		if err ==sql.ErrNoRows{
            http.NotFound(w,r) // verifica que el response y el request no sean nulos
			return
		}else{
			log.Fatal(w,err.Error(),http.StatusInternalServerError)
			return
		}
	 }
        filePath:=filepath.Join(uploadPath,fileName)//creamos el alojamiento de la imagen el cual recibe como parametros la carpeta y el nombre del archivo
   

		 //Verifica si el archivo existe dentro de la carpeta uploads
	    if _,err:=os.Stat(filePath);  os.IsNotExist(err){
			http.NotFound(w,r)
			return
		}

		// Lo guarda en el servidor de la imagen
	     http.ServeFile(w,r,filePath)
	
   
}

func postImagenCliente(w http.ResponseWriter,r *http.Request){
     
	err:=r.ParseMultipartForm(10<<20)//las imagenes tienen que ser menor a 10mb
	if err!=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	
	
	file,fileHeader,err:=r.FormFile("image")
	if err!=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
	defer file.Close()
	Id:=r.FormValue("Id")
	fileExt:=filepath.Ext(fileHeader.Filename)
	fileName:=fmt.Sprintf("%s%s",Id,fileExt)
	filePath:=filepath.Join("uploads",fileName)


	out,err:=os.Create(filePath)
	if err != nil {
		http.Error(w, "el archivo no se envio correctamente", http.StatusInternalServerError)
		return
	}
	 defer out.Close()
	_,err=io.Copy(out,file)
	 if err != nil {
		http.Error(w, "No se pudo salvar el archivo", http.StatusInternalServerError)
		return
	}
	_,err=db.Exec(`
	MERGE INTO imagenCliente AS target
	USING(SELECT @p1 AS Id, @p2 AS nombreImagen)
	AS source ON (target.Id  = source.Id)
	WHEN MATCHED THEN UPDATE SET nombreImagen=source.nombreImagen
	WHEN NOT MATCHED THEN INSERT (Id,nombreImagen) VALUES(source.Id, source.nombreImagen);`,Id,fileName)//permite si existe una fila con el mismo  parametro del id ingresado actualiza la imagen si no ingresa el elemento
	if err!=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Imagen Cliente Enviada correctamente"))
}