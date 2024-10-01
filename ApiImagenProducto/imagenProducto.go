package main

import(
	"database/sql"
	
	"net/http"
	"fmt"
	"log"
	"os"
	"io"
	"path/filepath"
_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)
  //inicializamos nuestro puntero de la base de datos
  var db *sql.DB
  const(
	uplodsPath="./uploads/"
  )
   //Creamos nuestra funcion para crear nuestra base de datos
   func initDB(){
     var err error
	 connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	 db,err=sql.Open("sqlserver",connectionString)
	 if err !=nil{
		log.Fatal("No se pudo conectar correctamente a la base de datos BizzBuz")

	 }

	 //Realizamos un pin para la activacion de la base de datos
	 err=db.Ping()
	 if err !=nil{
		log.Fatal(err)
	 }
	 fmt.Println("Se ha conectado correctamente a la base de Datos BizzBuz")
   }
   func main() {
	//llamamos funcion para inicializar la base de datos
	initDB();
	router:=mux.NewRouter();
	router.HandleFunc("/imagenProducto/{Id}",getImagenProducto).Methods("GET")
	router.HandleFunc("/imagenProducto",postImagenProducto).Methods("POST")
	headerOk:=handlers.AllowedHeaders([]string{"x-Requested-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
	log.Fatal(http.ListenAndServe(":5400",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

func getImagenProducto(w http.ResponseWriter, r * http.Request){
   vars:=mux.Vars(r)
   Id:=vars["Id"]
     var fileName string
   err:=db.QueryRow("SELECT nombreImagen FROM imagenProducto WHERE Id=@p1",Id).Scan(&fileName)
   if err !=nil{
	if err ==sql.ErrNoRows{
		http.NotFound(w,r)//verificamos que el cuerpo de la imagen no se encuentren nulos
		return
	}else{
		log.Fatal(w,err.Error(),http.StatusInternalServerError);
		return
	}
   }
   filePath:=filepath.Join(uplodsPath,fileName)
   if _,err:=os.Stat(filePath); os.IsNotExist(err){
         http.NotFound(w,r)
		 return;
   }
   http.ServeFile(w,r,filePath);
}

func postImagenProducto(w http.ResponseWriter, r *http.Request){
    err:=r.ParseMultipartForm(10<<20)//Verificamos que la imagen sea entre 10 y 20 mb
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
		return
	}

	file,fileHeader,err:=r.FormFile("image")
	 if err!=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
	    return
	}
	//cerramos conexion al momenti de encontrar la imagen
	defer file.Close();
	Id:=r.FormValue("Id")
	formExt:=filepath.Ext(fileHeader.Filename)
	fileName:=fmt.Sprintf("%s%s",Id,formExt)
	filePath:=filepath.Join("uploads",fileName)
	//Creamoa el path
	out,err:=os.Create(filePath);
     if err !=nil{
		log.Fatal(w,"Error al enviar la imagen",http.StatusInternalServerError)
		return
	 }
	 defer out.Close()
	 _,err=io.Copy(out,file)
	 if err!=nil{
      log.Fatal(w,"Error el formato de la imagen no es valida",http.StatusInternalServerError)
	  return
	 }

	 _,err=db.Exec(`
	   MERGE INTO imagenProducto AS target
	   USING(SELECT @p1 AS Id, @p2 AS nombreImagen)
	   AS source ON(target.Id=source.Id)
	   WHEN MATCHED THEN UPDATE SET nombreImagen=source.nombreImagen
	   WHEN NOT MATCHED THEN INSERT(Id,nombreImagen)VALUES(source.Id,source.nombreImagen);
	 `,Id,fileName)
	 if err !=nil{
		log.Fatal(w,err.Error(),http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("La imagen se envio Correctamente"))
}   