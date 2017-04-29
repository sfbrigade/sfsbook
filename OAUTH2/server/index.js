const express = require('express');
const cors = require('cors');

const app = express()
.use('cors')
.use('/auth', express.static('../client/dist/auth'));

app.listen(3000, () => console.log('server listening on port 3000'));