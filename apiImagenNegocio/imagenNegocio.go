package main
import(
	"database/sql"
	"net/http"
	"fmt"
	"log"
	"io"
	"os"
	"path/filepath"
_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)
//inicializamos nuestro puntero de nuestra conexion
 var db *sql.DB

 const(
	uploadsPath="./uploads/"
 )
 //creamos nuestra funcion para conectar base de datos
 func initDB(){
	var err error
	connectionString:="server=DESKTOP-58HBMDE;user=sserver;password=root;database=BizzBuz"
	db,err=sql.Open("sqlserver",connectionString)
	if err!=nil{
		log.Fatal("No se pudo conectar correctamente a la base de datos BizzBuz")
	}
	err=db.Ping()//Verificamos si la conexion esta activa
	if err !=nil{
		log.Fatal(err)
	}
	fmt.Println("Se ha conectado correctamente a la base de datos BizBuz...")
 }

func main() {
	//Llamamos nuestra cadena de conexion
	initDB()
	router:=mux.NewRouter()
	router.HandleFunc("/imagenNegocio/{Id}",getImagenNegocio).Methods("GET")
	router.HandleFunc("/imagenNegocio",postImagenNegocio).Methods("POST")
	headerOk:=handlers.AllowedHeaders([]string{"X-Requested-With","Content-Type","Authorization"})
	originOk:=handlers.AllowedOrigins([]string{"*"})
	methodsOk:=handlers.AllowedMethods([]string{"GET","HEAD","POST","PUT","OPTIONS"})
	log.Fatal(http.ListenAndServe(":5600",handlers.CORS(headerOk,originOk,methodsOk)(router)))
}

func getImagenNegocio(w http.ResponseWriter,r *http.Request){

	vars:=mux.Vars(r)
	Id:=vars["Id"]
	var fileName string
	err:= db.QueryRow("SELECT nombreImagen From imagenDescripcionNegocio WHERE Id=@p1",Id).Scan(&fileName)
	if err !=nil{
		if err ==sql.ErrNoRows{
			http.NotFound(w,r)//Verificamos que no hayan datos nulos o que no existan
			return
		}else{
			log.Fatal(w, err.Error(),http.StatusInternalServerError)
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

func postImagenNegocio(w http.ResponseWriter, r *http.Request){
  

	err:= r.ParseMultipartForm(10<<20)//verificamos de que la imagen este entre el rango de 10 y 20 mb
	if err !=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
		return
	}

	file,fileHeader,err:=r.FormFile("image")
	if err!=nil{
		log.Fatal(w,err.Error(),http.StatusBadRequest)
		return
	}
	defer file.Close()
	Id:=r.FormValue("Id")
	formExt:=filepath.Ext(fileHeader.Filename)
	fileName:=fmt.Sprintf("%s%s",Id,formExt)
	filePath:=filepath.Join("uploads",fileName)

	//Creamos el path
	out,err:=os.Create(filePath)
	if err !=nil{
		log.Fatal(w,"error al enviar la imagen formato no permitido",http.StatusInternalServerError)
		return
	}
	defer out.Close();

	_,err=io.Copy(out,file)
	if err !=nil{
		log.Fatal(w,"Imagen invalida no es el archivo correspondiente",http.StatusInternalServerError)
		return
	}
	_,err=db.Exec(`
	MERGE INTO imagenDescripcionNegocio AS target
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