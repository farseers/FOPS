import{k as _,ah as t,o,c as a,W as r,Q as l,P as c,F as C,a9 as v,a as h,K as w,T as i}from"./@vue-d9027515.js";import{_ as S}from"./_plugin-vue_export-helper-c27b6911.js";const B={class:"icon-selector-warp-row"},b=_({name:"iconSelectorList"}),I=_({...b,props:{list:{type:Array,default:()=>[]},empty:{type:String,default:()=>"无相关图标"},prefix:{type:String,default:()=>""}},emits:["get-icon"],setup(e,{emit:m}){const p=e,d=s=>{m("get-icon",s)};return(s,N)=>{const u=t("SvgIcon"),f=t("el-col"),g=t("el-row"),y=t("el-empty"),k=t("el-scrollbar");return o(),a("div",B,[r(k,{ref:"selectorScrollbarRef"},{default:l(()=>[p.list.length>0?(o(),c(g,{key:0,gutter:10},{default:l(()=>[(o(!0),a(C,null,v(e.list,(n,x)=>(o(),c(f,{xs:6,sm:4,md:4,lg:4,xl:4,key:x,onClick:V=>d(n)},{default:l(()=>[h("div",{class:w(["icon-selector-warp-item",{"icon-selector-active":e.prefix===n}])},[r(u,{name:n},null,8,["name"])],2)]),_:2},1032,["onClick"]))),128))]),_:1})):i("",!0),e.list.length<=0?(o(),c(y,{key:1,"image-size":100,description:e.empty},null,8,["description"])):i("",!0)]),_:1},512)])}}});const L=S(I,[["__scopeId","data-v-e4865837"]]);export{L as default};