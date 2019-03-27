import React from 'react';

let GlobalContext;

export default function GetGlobalContext() {
  if (!GlobalContext) {
    GlobalContext = React.createContext();
  }
  return GlobalContext;
}
