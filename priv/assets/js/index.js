new Vue({
  el: "#app",
  data: {
    title: "Authentikos",
    qrCode: "",
    flags: {
      passcode: false,
      checkResult: false
    },
    generateInputs: {
      account_name: ""
    },
    verifyInputs: {
      account_name: "",
      passcode    : ""
    }
  },
  methods: {
    generateQrCode: function(e) {
      e.preventDefault();
      var params = this.generateInputs;
      this.$http.get("/otp/generate", { params: params })
        .then((response) => {
          qrBlob = response.body;
          this.flags.checkResult = false;
          this.qrCode = URL.createObjectURL(qrBlob);
        }).catch((error) => {
          console.log(error);
        });
    },

    verifyQrCode: function(e) {
      e.preventDefault();
      var params = this.verifyInputs;
      this.$http.get("/otp/verify", { params: params })
        .then((response) => {
          return response.json();
        }).catch((error) => {
          console.log(error);
        }).then((json) => {
          status = json.code;
          this.flags.checkResult = true;
          if (status != "OK") {
            this.flags.passcode = false;
          } else {
            this.flags.passcode = true;
          }
        })
    },

    clearOutput: function(e, type) {
      if (type === "generate") {
        this.qrCode = "";
      } else if (type === "verify") {
        this.flags.checkResult = false;
      }
    }

  }
})
