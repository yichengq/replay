(function() {
  function playground(opts) {
    var PROMPT = "> " ;
    var log_output = PROMPT;
    var ERR_OUTPUT = "Warnings: "
    var err_output = ERR_OUTPUT;
    var original_log = console.log;
    var original_err = console.error;
    function reset_log_output (){ log_output  = PROMPT;}
    function reset_error_output(){ err_output = ERR_OUTPUT;}
    function redirect_err() {
      err_output = err_output + Array.prototype.slice.apply(arguments).join(' ') + "\n"
    };
    function redirect() { log_output = log_output + Array.prototype.slice.apply(arguments).join(' ') + "\n"};
    function get_log_output(){
      var old = log_output;
      reset_log_output();
      return old;
    }
    function get_error_output(){
      var old = err_output;
      reset_error_output();
      return old;
    }
    function evalCode(js){
      console.log = redirect;
      try {
        window.eval(js);
        console.log = original_log;
        opts.outputMirror.setValue(get_log_output());
      }
      catch(e){
        opts.outputMirror.setValue(get_log_output() + "\n" + e);
        console.log = original_log;
      }
    }

    function origin(href) {
      return (""+href).split("/").slice(0, 3).join("/");
    }

    function body() {
      return opts.codeMirror.getValue();
    }

    function setBody(text) {
      opts.codeMirror.setValue(text);
    }

    function run() {
      $.ajax("/compile", {
        data: body(),
        type: "POST",
        dataType: "text",
        complete: function(xhr) {
          if (xhr.status == 200) {
            evalCode(xhr.responseText);
            return;
          }
          if (xhr.status == 400) {
            opts.outputMirror.setValue(xhr.responseText);
            return;
          }
          alert("Server error; try again.\nError: " + xhr.responseText);
        }
      });
    }

    function fmt() {
      $.ajax("/format", {
        data: body(),
        type: "POST",
        dataType: "text",
        complete: function(xhr) {
          if (xhr.status == 200) {
            setBody(xhr.responseText);
            opts.outputMirror.setValue("");
            return;
          }
          if (xhr.status == 400) {
            opts.outputMirror.setValue(xhr.responseText);
            return;
          }
          alert("Server error; try again.\nError: " + xhr.responseText);
        }
      });
    }

    var shareURL = $(opts.shareURLEl).hide();
    function share() {
      var sharingData = body();
      $.ajax("/share", {
        processData: false,
        data: sharingData,
        type: "POST",
        complete: function(xhr) {
          sharing = false;
          if (xhr.status != 200) {
            alert("Server error; try again.");
            return;
          }
          if (shareURL) {
            var path = "/p/" + xhr.responseText;
            var url = origin(window.location) + path;
            shareURL.show().val(url).focus().select();
            var historyData = {"code": sharingData};
            window.history.pushState(historyData, "", path);
            pushedEmpty = false;
          }
        }
      });
    }

    var pushedEmpty = (window.location.pathname == "/");
    function onInputChanges() {
      if (pushedEmpty) {
        return;
      }
      pushedEmpty = true;
      $(opts.shareURLEl).hide();
      window.history.pushState(null, "", "/");
    }


    opts.codeMirror.on("changes", onInputChanges);
    $(opts.shareEl).click(share);
    $(opts.runEl).click(run);
    $(opts.fmtEl).click(fmt);
  }

  window.playground = playground;
})();
