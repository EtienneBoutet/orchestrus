<template>
  <v-dialog v-model="showSelf" max-width="300px">
    <template v-slot:activator="{ on, attrs }">
      <v-btn
        icon
        color="primary"
        v-bind="attrs"
        v-on="on"
      >
        <v-icon dark>mdi-plus</v-icon>
      </v-btn>
    </template>
    <v-card>
      <v-card-title>
        <span class="headline">Add an image</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field label="Name" v-model="name" required></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field label="Public Port" v-model="publicPort" required></v-text-field>
            </v-col>
            <v-col cols="6">
              <v-text-field label="Local Port" v-model="localPort" required></v-text-field>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="showSelf = false">Close</v-btn>
        <v-btn color="blue darken-1" text @click="addImage">Add</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  name: "AddImage",
  props: {
    worker: Object
  },
  data () {
    return {
      showSelf: false,
      name: "",
      publicPort: "",
      localPort: ""
    }
  },
  methods: {
    addImage: function () {
      
      // Do stuff about adding.
      // Create axios request

      let id = ""

      this.$emit("added", {worker: this.worker, name: this.name, id: id, publicPort: this.publicPort, localPort: this.localPort})

      this.name = ""
      this.publicPort = ""
      this.localPort = ""
      this.showSelf = false
    }
  }
}
</script>