<template>
  <div class="container">
    <div class="images-header">
      <div>
        <h1>Orchestrus</h1>
      </div>
      <div>
        <v-btn icon color="primary" v-on:click="refreshWorkers">
          <v-icon>mdi-cached</v-icon>
        </v-btn>
      </div>
    </div>
    <div class="images-header">
      <div>
        <h2>Workers</h2>
      </div>
      <div>
        <AddWorker @added="addWorker"/>
      </div>
    </div>
    <div class="images-header">
      <v-data-table 
        :hide-default-footer="true"
        :headers="workersHeaders"
        :items="workers"
        item-key="ip"
        show-expand
        class="elevation-1"
      >
        <template v-slot:[`item.status`]="{ item }">
          <v-simple-checkbox v-model="item.status" disabled></v-simple-checkbox>
        </template>
        <template v-slot:expanded-item="{ item }">
          <div class="images-header">
            <div>
              <h3>Images</h3>
            </div>
            <div>
              <AddImage :worker="item" @added="addImage"/>
            </div>
          </div>
          <div class="images-header">
            <v-data-table 
              :hide-default-footer="true"
              :headers="imagesHeaders"
              :items="item.images"
              class="elevation-1"
            >
              <template v-slot:item="row">
                <tr>
                  <td>{{row.item.name}}</td>
                  <td>{{row.item.id}}</td>
                  <td>{{row.item.port}}</td>
                  <td>
                    <v-btn class="mx-2" fab dark small color="red" @click="removeImage(row.item, item)">
                      <v-icon dark>mdi-delete</v-icon>
                    </v-btn>
                  </td>
                </tr>
              </template>
            </v-data-table>
          </div>
        </template>
      </v-data-table>
    </div>
  </div>
</template>

<script>

import AddWorker from "./AddWorker";
import AddImage from "./AddImage";

export default {
  name: 'Application',
  components: {
    AddWorker,
    AddImage
  },
  data () {
    return {
      expanded: [],
      workersHeaders: [
        { text: 'IP', align: 'start', value: 'ip'},
        { text: 'Status', value: 'status' },
      ],
      imagesHeaders: [
        { text: 'Name', align: 'start', value: 'name'},
        { text: 'ID', value: 'id'},
        { text: 'Port', value: 'port'}
      ],
      workers: [
        {
          ip: "127.0.0.1",
          status: true,
          images: [
            {
              name: "httpd",
              id: "466c47d39c33acd3b4ba0281c9bc145de1e8744dee12d728fb39b67f906bf4d4",
              port: "8081:80"
            }
          ]
        },
        {
          ip: "192.168.0.1",
          status: false
        }
      ]
    }
  },
  methods: {
    refreshWorkers: function () {
      alert("Refreshing workers")
    },
    addWorker: function (workerIp) {
      this.workers.push({
          ip: workerIp,
          images: []
      })
    },

    addImage: function (imageInfo) {
      const worker = imageInfo.worker

      worker.images.push({
        name: imageInfo.name,
        id: imageInfo.id,
        port: imageInfo.publicPort.concat(":", imageInfo.localPort)
      })
    },

    removeImage: function (image, worker) {
      worker.images = worker.images.filter(function( obj ) {
        return obj !== image;
      });
    },

  }


}
</script>