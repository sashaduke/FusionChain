import { Outlet } from "react-router-dom";
import Menu from "./menu";
import { useEffect, useState } from "react";
import { enableKeplr, useKeplrAddress, useMetamaskAddress } from "../keplr";
import AccountInfo from "@/components/account_info";
import { Toaster } from "@/components/ui/toaster";

function sleep(ms: number) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function Root() {
  useEffect(() => {
    window.onload = async () => {
      // await enableKeplr();
    }
  }, []);

  return (
    <div>
      <header className="flex flex-row justify-between py-2 px-4 border-b border-slate-200">
        <Menu />
        <AccountInfo />
      </header>
      <Outlet />
      <Toaster />
    </div>
  );
}

export default Root;
