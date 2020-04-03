axios.defaults.baseURL = '';
axios.defaults.headers.post['Content-Type'] = 'application/x-www-form-urlencoded';
axios.defaults.headers.common['Token'] = localStorage.getItem('cycron-token');


function getUrlKey(name,url){
    return decodeURIComponent((new RegExp('[?|&]' + name + '=' + '([^&;]+?)(&|#|;|$)').exec(url) || [, ""])[1].replace(/\+/g, '%20')) || null
}
function trimStr(str){
    if(typeof(str) != 'string'){
        return ''
    }

    return str.replace(/(^\s*)|(\s*$)/g,"")
}
