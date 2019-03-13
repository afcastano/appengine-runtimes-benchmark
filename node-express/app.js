'use strict';

const express = require('express');

const app = express();
app.enable('trust proxy');

// By default, the client will authenticate using the service account file
// specified by the GOOGLE_APPLICATION_CREDENTIALS environment variable and use
// the project specified by the GOOGLE_CLOUD_PROJECT environment variable. See
// https://github.com/GoogleCloudPlatform/google-cloud-node/blob/master/docs/authentication.md
// These environment variables are set automatically on Google App Engine
const Datastore = require('@google-cloud/datastore');

// Instantiate a datastore client
const datastore = Datastore();

function getEntities(random2) {
  const query = datastore
    .createQuery('DummyEntity')
    .filter('random2', '>=', random2)
    .filter('random2', '<', random2 + 10000)
    .limit(10);

  return datastore.runQuery(query);
}

app.get('/', async (req, res) => {
  console.log('Checking service running');
  res
      .status(200)
      .set('Content-Type', 'text/plain')
      .send(`Express server running v4`)
      .end();  
});

app.get('/entities/:random2', async (req, res, next) => {
  console.log('Loading entities');
  try {
    const results = await getEntities(parseInt(req.params.random2));
    const entities = results[0].map(entity => {
      const id = entity[datastore.KEY].path[1];
      return {
        id,
        ...entity
      }
    });
    
    res
      .status(200)
      .set('Content-Type', 'application/json')
      .send(JSON.stringify(entities))
      .end();
  } catch (error) {
    console.error(`ERROR ${error}`);
    next(error);
  }
});

const PORT = process.env.PORT || 8080;
app.listen(process.env.PORT || 8080, () => {
  console.log(`App listening on port ${PORT}`);
  console.log('Press Ctrl+C to quit.');
});

module.exports = app;
