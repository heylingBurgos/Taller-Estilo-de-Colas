const express = require("express");
const cors = require("cors");
const axios = require("axios");
const path = require('path');

const app = express();

/*var corsOptions1 = {
    origin: "http://localhost:8081/"
};

var corsOptions2 = {
    origin: "http://localhost:9092"
};*/

app.use(cors({ origin: 'http://localhost:8080' }));

// parse requests of content-type - application/json
app.use(express.json());

// parse requests of content-type - application/x-www-form-urlencoded
app.use(express.urlencoded({ extended: true }));

// simple route
app.get("/", (req, res) => {
    res.json({ message: "Servidor ejecutandose correctamente" });
});


// recibir datos desde cliente componente DigiturnoUser
app.post("/",(req, res) => {
    
    const datos = {
        id: req.body.id,
        name: req.body.name,
        cellphone: req.body.cellphone
    };
    console.log(datos);
    //Enviar datos a backend
    axios.post("http://localhost:5000", datos)
    .then(response => {
        console.log(response.data);
        res.json({mensaje: "Datos recibidos del backend correctamente"})
    })
    .catch(error => {
        console.log(error);
        res.json({mensaje: "Error al enviar datos"})
    })

    res.json({mensaje: "Datos recibidos correctamente"})
})

const history = require('connect-history-api-fallback');
app.use(history());
app.use(express.static(path.join(__dirname, 'public')));


// set port, listen for requests
const PORT = process.env.PORT || 4000;
app.listen(PORT, () => {
    console.log(`Server is running on port ${PORT}.`);
});