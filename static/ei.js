(function(){
  'use strict';
  function enjectAll() {
    for (var embeds=document.getElementsByClassName("replay"), i=0; i<embeds.length; i++){
      var data = getDataAttributes(embeds[i]);
      console.log(data);
      var content = createEmbededIframe(data);
      enjectElement(embeds[i],content);
    }
  }

  function getDataAttributes(e) {
    var output = {};
    for(var attributes=e.attributes, i=0;i<attributes.length;i++){
      var name = attributes[i].name;
      if (name.indexOf("data-")===0){
      console.log(name);
        output[name.replace("data-","")]=attributes[i].value;
      }
    }
    return output;
  }

  function createEmbededIframe(data) {
    var attributes = {
      id: "re_embed_"+data.id,
      src: protocol()+"//"+data.host+"/p/"+data.id,
      scrolling:"no",
      frameborder:"0",
      name:"Reaon Playground Embed",
      title: "",
      allowTransparency:"true",
      allowfullscreen:"true",
      width: "100%",
      height: data.height?data.height:"300",
      class: "re_embed_iframe",
      style: "overflow: hidden;"
    };
    var output = "<iframe ";
    for (var k in attributes) {
      output += k+'="'+attributes[k]+'" ';
    }
    return output+"></iframe>";
  }

  function protocol() {
    var protocol = document.location.protocol;
    if (protocol==="file:") {
      protocol = "http:";
    }
    return protocol;
  }

  function enjectElement(curDom, content) {
    if (curDom.parentNode) {
      var wrapper = document.createElement("div");
      wrapper.className = "re_embed_wrapper";
      wrapper.innerHTML = content;
      curDom.parentNode.replaceChild(wrapper, curDom);
      return;
    }
    curDom.innerHTML = content;
  }

  function docReadyRun(f) {
    /in/.test(document.readyState)?setTimeout("window.__re_docRun("+f+")",9):f();
  }

  window.__re_docRun = docReadyRun;
  window.__re_enjectAll = enjectAll;
  docReadyRun(function(){window.__re_enjectAll();});
})();
