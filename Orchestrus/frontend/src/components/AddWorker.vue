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
        <span class="headline">Add a worker</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12">
              <v-text-field label="IP" v-model="ip" required></v-text-field>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="showSelf = false">Close</v-btn>
        <v-btn color="blue darken-1" text @click="addWorker">Save</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  name: "AddWorker",
  data () {
    return {
      showSelf: false,
      ip: "",
    }
  },
  methods: {
    addWorker: function () {
      
      // Do stuff about adding.
      // Create axios request


      this.$emit("added", this.ip)

      this.ip = ""
      this.showSelf = false
    }
  }
}
</script>