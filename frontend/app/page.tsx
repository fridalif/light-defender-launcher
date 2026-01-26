"use client"

import type React from "react"

import { useEffect, useState } from "react"
import Image from "next/image"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent } from "@/components/ui/card"
import { Spinner } from "@/components/ui/spinner"
import { useLanguage } from "@/contexts/language-context"
import { AlertCircle, Eye, EyeOff, Upload, ChevronDown, ChevronUp, PlayCircle } from "lucide-react"
import { LanguageSelector } from "@/components/language-selector"
import { InstallLightDefender, LoadConfig, UpdateLightDefender } from "@/wailsjs/go/main/App"
import { EventsOn } from "@/wailsjs/runtime/runtime"

type ExecutionMode = "install" | "reconfigure" | "update"
type ConfigMethod = "file" | "credentials"
type ExecutionStatus = "idle" | "running" | "success" | "failure"

export default function LauncherPage() {
  const { t } = useLanguage()

  // Authentication
  const [sshIp, setSshIp] = useState("")
  const [sshLogin, setSshLogin] = useState("root")
  const [sshPassword, setSshPassword] = useState("")
  const [showSshPassword, setShowSshPassword] = useState(false)
  const [customPort, setCustomPort] = useState(false)
  const [sshPort, setSshPort] = useState("22")

  // Mode selection
  const [mode, setMode] = useState<ExecutionMode>("install")

  // Config method (for install and reconfigure modes)
  const [configMethod, setConfigMethod] = useState<ConfigMethod>("credentials")
  const [configFile, setConfigFile] = useState<File | null>(null)
  const [dashboardLogin, setDashboardLogin] = useState("")
  const [dashboardPassword, setDashboardPassword] = useState("")
  const [showDashboardPassword, setShowDashboardPassword] = useState(false)
  const [configId, setConfigId] = useState("")

  // Execution status
  const [status, setStatus] = useState<ExecutionStatus>("idle")
  const [showDetails, setShowDetails] = useState(false)
  const [executionLog, setExecutionLog] = useState<string[]>([])

  const handleFileUpload = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setConfigFile(e.target.files[0])
    }
  }

  useEffect(() => {
    EventsOn("log", (data: string) => {
      setExecutionLog((log) => [...log, data])
    })
  },[])

  const handleStart = async () => {
    setStatus("running")
    setShowDetails(true)
    setExecutionLog([])
    if (mode == "update") {
      let data = await UpdateLightDefender(sshLogin, sshPassword, sshIp, sshPort)
      if (data[0] == "true") {
        setStatus("success")
      } else {
        setStatus("failure")
        setExecutionLog((log) => [...log, data[1]])
      }
    }
    if (mode == "reconfigure") {
      let data = await LoadConfig(sshLogin, sshPassword, sshIp, sshPort, dashboardLogin, dashboardPassword, configId, configFile != null ? (await configFile.bytes()).toBase64() : "")
      if (data[0] == "true") {
        setStatus("success")
      } else {
        setStatus("failure")
        setExecutionLog((log) => [...log, data[1]])
      }
    }
    if (mode == "install") {
      let data = await InstallLightDefender(sshLogin, sshPassword, sshIp, sshPort)
      if (data[0] == "false") {
        setStatus("failure")
        setExecutionLog((log) => [...log, data[1]])
        return
      }
      let data2 = await LoadConfig(sshLogin, sshPassword, sshIp, sshPort, dashboardLogin, dashboardPassword, configId, configFile != null ? (await configFile.bytes()).toBase64() : "")
      if (data2[0] == "true") {
        setStatus("success")
      } else {
        setStatus("failure")
        setExecutionLog((log) => [...log, data2[1]])
      }
    }
  }

  const getStatusColor = () => {
    switch (status) {
      case "idle":
        return "text-slate-400"
      case "running":
        return "text-blue-400"
      case "success":
        return "text-green-400"
      case "failure":
        return "text-red-400"
    }
  }

  const getStatusText = () => {
    switch (status) {
      case "idle":
        return t("launcher.status.idle")
      case "running":
        return t("launcher.status.running")
      case "success":
        return t("launcher.status.success")
      case "failure":
        return t("launcher.status.failure")
    }
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-950 via-slate-900 to-slate-950 p-4 md:p-6 lg:p-8">
      {/* Language selector */}
      <div className="absolute top-4 right-4 md:top-6 md:right-6">
        <LanguageSelector />
      </div>

      <div className="max-w-4xl mx-auto space-y-6">
        {/* Logo */}
        <div className="flex flex-col items-center space-y-4 pt-8 md:pt-12">
          <div className="relative w-20 h-20 md:w-24 md:h-24">
            <Image src="/images/1000001892-cropped.png" alt="Light Defender Logo" fill className="object-contain" />
          </div>
          <h1 className="text-2xl md:text-3xl font-bold text-white">Light Defender Launcher</h1>
        </div>

        {/* Main Card */}
        <Card className="border-slate-800 bg-slate-900/50 backdrop-blur-sm shadow-2xl">
          <CardContent className="p-4 md:p-6 space-y-6">
            {/* SSH Credentials */}
            <div className="space-y-4">
              <div className="flex items-center gap-2 text-sm text-slate-400">
                <AlertCircle className="h-4 w-4" />
                <span>{t("launcher.localPasswordNotice")}</span>
              </div>

              {/* SSH IP field */}
              <div className="space-y-2">
                <Label htmlFor="ssh-ip" className="text-slate-200">
                  {t("launcher.sshIp")}
                </Label>
                <Input
                  id="ssh-ip"
                  type="text"
                  value={sshIp}
                  onChange={(e) => setSshIp(e.target.value)}
                  placeholder={t("launcher.sshIpPlaceholder")}
                  className="bg-slate-800/50 border-slate-700 text-white"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <Label htmlFor="ssh-login" className="text-slate-200">
                    {t("launcher.sshLogin")}
                  </Label>
                  <Input
                    id="ssh-login"
                    type="text"
                    value={sshLogin}
                    onChange={(e) => setSshLogin(e.target.value)}
                    className="bg-slate-800/50 border-slate-700 text-white"
                  />
                </div>

                <div className="space-y-2">
                  <Label htmlFor="ssh-password" className="text-slate-200">
                    {t("launcher.sshPassword")}
                  </Label>
                  <div className="relative">
                    <Input
                      id="ssh-password"
                      type={showSshPassword ? "text" : "password"}
                      value={sshPassword}
                      onChange={(e) => setSshPassword(e.target.value)}
                      className="bg-slate-800/50 border-slate-700 text-white pr-10"
                    />
                    <button
                      type="button"
                      onClick={() => setShowSshPassword(!showSshPassword)}
                      className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200"
                    >
                      {showSshPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                    </button>
                  </div>
                </div>
              </div>

              {/* Custom port checkbox and port field */}
              <div className="space-y-3">
                <div className="flex items-center space-x-2">
                  <input
                    type="checkbox"
                    id="custom-port"
                    checked={customPort}
                    onChange={(e) => setCustomPort(e.target.checked)}
                    className="w-4 h-4 rounded border-slate-700 bg-slate-800/50 text-blue-600 focus:ring-2 focus:ring-blue-500 focus:ring-offset-0"
                  />
                  <Label htmlFor="custom-port" className="text-slate-200 cursor-pointer">
                    {t("launcher.customPort")}
                  </Label>
                </div>

                {customPort && (
                  <div className="space-y-2">
                    <Label htmlFor="ssh-port" className="text-slate-200">
                      {t("launcher.sshPort")}
                    </Label>
                    <Input
                      id="ssh-port"
                      type="text"
                      value={sshPort}
                      onChange={(e) => setSshPort(e.target.value)}
                      placeholder="22"
                      className="bg-slate-800/50 border-slate-700 text-white"
                    />
                  </div>
                )}
              </div>
            </div>

            {/* Mode Selection */}
            <div className="space-y-3">
              <Label className="text-slate-200">{t("launcher.mode")}</Label>
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                <button
                  onClick={() => setMode("install")}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    mode === "install"
                      ? "border-blue-500 bg-blue-500/10 text-blue-400"
                      : "border-slate-700 bg-slate-800/30 text-slate-400 hover:border-slate-600"
                  }`}
                >
                  <div className="font-medium">{t("launcher.mode.install")}</div>
                </button>
                <button
                  onClick={() => setMode("reconfigure")}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    mode === "reconfigure"
                      ? "border-blue-500 bg-blue-500/10 text-blue-400"
                      : "border-slate-700 bg-slate-800/30 text-slate-400 hover:border-slate-600"
                  }`}
                >
                  <div className="font-medium">{t("launcher.mode.reconfigure")}</div>
                </button>
                <button
                  onClick={() => setMode("update")}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    mode === "update"
                      ? "border-blue-500 bg-blue-500/10 text-blue-400"
                      : "border-slate-700 bg-slate-800/30 text-slate-400 hover:border-slate-600"
                  }`}
                >
                  <div className="font-medium">{t("launcher.mode.update")}</div>
                </button>
              </div>
            </div>

            {/* Config Loading Source (only for install and reconfigure) */}
            {(mode === "install" || mode === "reconfigure") && (
              <div className="space-y-4">
                <Label className="text-slate-200">{t("launcher.loadingSource")}</Label>

                {/* Method selector */}
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <button
                    onClick={() => setConfigMethod("file")}
                    className={`p-3 rounded-lg border-2 transition-all ${
                      configMethod === "file"
                        ? "border-blue-500 bg-blue-500/10 text-blue-400"
                        : "border-slate-700 bg-slate-800/30 text-slate-400 hover:border-slate-600"
                    }`}
                  >
                    <Upload className="h-5 w-5 mx-auto mb-1" />
                    <div className="text-sm font-medium">{t("launcher.loadingSource.uploadFile")}</div>
                  </button>
                  <button
                    onClick={() => setConfigMethod("credentials")}
                    className={`p-3 rounded-lg border-2 transition-all ${
                      configMethod === "credentials"
                        ? "border-blue-500 bg-blue-500/10 text-blue-400"
                        : "border-slate-700 bg-slate-800/30 text-slate-400 hover:border-slate-600"
                    }`}
                  >
                    <div className="text-sm font-medium">{t("launcher.loadingSource.downloadFromServer")}</div>
                  </button>
                </div>

                {/* File upload */}
                {configMethod === "file" && (
                  <div className="space-y-2">
                    <Label htmlFor="config-file" className="text-slate-200">
                      {t("launcher.configFile")}
                    </Label>
                    <div className="flex items-center gap-3">
                      <Input
                        id="config-file"
                        type="file"
                        onChange={handleFileUpload}
                        className="bg-slate-800/50 border-slate-700 text-white file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:bg-blue-600 file:text-white file:cursor-pointer hover:file:bg-blue-700 cursor-pointer"
                      />
                    </div>
                    {configFile && (
                      <p className="text-sm text-slate-400">
                        {t("launcher.selectedFile")}: {configFile.name}
                      </p>
                    )}
                  </div>
                )}

                {/* Dashboard credentials */}
                {configMethod === "credentials" && (
                  <div className="space-y-4">
                    <div className="space-y-2">
                      <Label htmlFor="dashboard-login" className="text-slate-200">
                        {t("launcher.dashboardLogin")}
                      </Label>
                      <Input
                        id="dashboard-login"
                        type="text"
                        value={dashboardLogin}
                        onChange={(e) => setDashboardLogin(e.target.value)}
                        placeholder={t("launcher.dashboardLoginPlaceholder")}
                        className="bg-slate-800/50 border-slate-700 text-white"
                      />
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="dashboard-password" className="text-slate-200">
                        {t("launcher.dashboardPassword")}
                      </Label>
                      <div className="relative">
                        <Input
                          id="dashboard-password"
                          type={showDashboardPassword ? "text" : "password"}
                          value={dashboardPassword}
                          onChange={(e) => setDashboardPassword(e.target.value)}
                          placeholder={t("launcher.dashboardPasswordPlaceholder")}
                          className="bg-slate-800/50 border-slate-700 text-white pr-10"
                        />
                        <button
                          type="button"
                          onClick={() => setShowDashboardPassword(!showDashboardPassword)}
                          className="absolute right-3 top-1/2 -translate-y-1/2 text-slate-400 hover:text-slate-200"
                        >
                          {showDashboardPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                        </button>
                      </div>
                    </div>

                    <div className="space-y-2">
                      <Label htmlFor="config-id" className="text-slate-200">
                        {t("launcher.configId")}
                      </Label>
                      <Input
                        id="config-id"
                        type="text"
                        value={configId}
                        onChange={(e) => setConfigId(e.target.value)}
                        placeholder={t("launcher.configIdPlaceholder")}
                        className="bg-slate-800/50 border-slate-700 text-white"
                      />
                    </div>
                  </div>
                )}
              </div>
            )}

            {/* Start Button */}
            <div className="flex justify-center pt-4">
              <button onClick={handleStart} disabled={status === "running"} className="relative group">
                <div
                  className={`w-32 h-32 md:w-40 md:h-40 rounded-full bg-gradient-to-br from-blue-600 to-blue-800 flex items-center justify-center transition-transform ${status === "running" ? "scale-95" : "group-hover:scale-105"} ${status === "running" ? "opacity-50" : ""}`}
                >
                  <div className="relative w-20 h-20 md:w-24 md:h-24">
                    {status === "running" ? (
                      <Spinner className="w-full h-full text-white" />
                    ) : (
                      <PlayCircle className="w-full h-full text-white" />
                    )}
                  </div>
                </div>
                <div className="absolute inset-0 rounded-full bg-blue-500/20 blur-xl group-hover:blur-2xl transition-all" />
              </button>
            </div>

            {/* Status */}
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <span className="text-slate-400 text-sm">{t("launcher.status.label")}: </span>
                  <span className={`font-medium ${getStatusColor()}`}>{getStatusText()}</span>
                </div>
                {status !== "idle" && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={() => setShowDetails(!showDetails)}
                    className="text-slate-400 hover:text-slate-200"
                  >
                    {showDetails ? (
                      <>
                        <ChevronUp className="h-4 w-4 mr-1" />
                        {t("launcher.hideDetails")}
                      </>
                    ) : (
                      <>
                        <ChevronDown className="h-4 w-4 mr-1" />
                        {t("launcher.showDetails")}
                      </>
                    )}
                  </Button>
                )}
              </div>

              {/* Execution log */}
              {showDetails && status !== "idle" && (
                <div className="bg-black/40 border border-slate-800 rounded-lg p-4 max-h-64 overflow-y-auto font-mono text-sm">
                  {executionLog.map((log, index) => (
                    <div key={index} className="text-green-400 mb-1">
                      {log}
                    </div>
                  ))}
                  {status === "running" && <div className="text-green-400 animate-pulse">▊</div>}
                </div>
              )}
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}
