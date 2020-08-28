<template>
  <div class="images-header">
    <div>
      <h3 class="title is-3">Workers</h3>
    </div>
    <div>
      <button>
        <span class="icon">
          <i class="fas fa-plus-circle"></i>
        </span>
      </button>
    </div>
  </div>
  <div class="images-header">
    <v-simple-table>
      <template v-slot:default>
        <thead>
          <tr>
            <th class="text-left">Name</th>
            <th class="text-left">Running images</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="host in hosts" :key="host.name">
            <td>{{ host.ip }}</td>
            <td>{{ host.images}}</td>
          </tr>
        </tbody>
      </template>
    </v-simple-table>
  </div>
</template>

<script>

export default {
  name: 'Workers',
   data () {
      return {
        hosts: [
          {
            ip: 'xxx.xxx.xxx.xxx',
            images: [
              "httpd",
              "apache"
            ]
          },
          {
            ip: 'yyy.yyy.yyy.yyy',
            images: [
              "httpd",
            ]
          }
        ],
      }
    },
}
</script>