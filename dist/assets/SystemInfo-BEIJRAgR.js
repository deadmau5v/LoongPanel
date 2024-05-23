/*
 * 创建人： deadmau5v
 * 创建时间： 2024-0-0
 * 文件作用：
 */

import{C as _,j as s}from"./index-tED-CDpK.js";import{r as i}from"./react-dom-D8wjeNdT.js";import{C as r}from"./index-DIeflrbH.js";import{D as e}from"./index-C7xbLzKn.js";import"./App-CzhtfASx.js";function d(){const n=i.useContext(_),a=`${n==null?void 0:n.API_URL}/api/v1/status/system_info`,[t,o]=i.useState({system_arch:"",system_public_ip:"",system_cpu_name:"",system_linux_version:"",system_run_time:""}),c=async()=>{const m=await fetch(a);if(m.ok){const l=await m.json();o(l)}};return i.useEffect(()=>{c()},[]),s.jsx(s.Fragment,{children:s.jsx(r,{id:"SystemInfo",className:"card",title:"系统信息",children:s.jsx(r.Meta,{description:s.jsxs(e,{size:"small",column:1,children:[s.jsx(e.Item,{label:"公网 IP",children:t.system_public_ip}),s.jsx(e.Item,{label:"系统架构",children:t.system_arch}),s.jsx(e.Item,{label:"处理器名",children:t.system_cpu_name}),s.jsx(e.Item,{label:"内核版本",children:t.system_linux_version}),s.jsx(e.Item,{label:"运行时间",children:t.system_run_time})]})})})})}export{d as default};
