new Vue({
  el: "#app",
  data: {
    title: "Authentikos",
    identity: "",
    qrCode: ""
  },
  methods: {
    generateQrCode: function(e) {
      e.preventDefault();
      var params = { account_name: this.identity }
      this.$http.get("/otp/generate", { params: params })
        .then((response) => {
          qrBlob = response.body;
          this.qrCode = URL.createObjectURL(qrBlob);
        });
    }
  }
})
