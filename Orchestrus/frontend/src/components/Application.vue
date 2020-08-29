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
        :headers="headers"
        :items="workers"
        class="elevation-1"
      ></v-data-table>
    </div>
  </div>
</template>

<script>

import AddWorker from "./AddWorker";

export default {
  name: 'Application',
  components: {
    AddWorker
  },
  data () {
    return {
      headers: [
        {
          text: 'IP',
          align: 'start',
          value: 'ip',
        },
        { text: 'Images', value: 'images' },
      ],
      workers: [
        {
          ip: "127.0.0.1",
          images: ["httpd, apache"]
        },
        {
          ip: "192.168.0.1",
          images: ["python"]
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
    }
  }


}
</script>