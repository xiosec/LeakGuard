// dllmain.cpp : Defines the entry point for the DLL application.
#include "pch.h"
#include <Windows.h>
#include <process.h>
#include <stdlib.h>
#include <stdio.h>
#include <SubAuth.h>
#include <wininet.h>
#include <codecvt>

#pragma comment(lib, "wininet")
#define MAX_PASSWORD_LENGTH 1024

using namespace std;

BOOL GetQueryValue(LPCSTR name, char* buffer, DWORD bufferSize) {
	HKEY hKey;
	DWORD dwType = REG_SZ;

	if (RegOpenKeyExA(HKEY_LOCAL_MACHINE, "SYSTEM\\CurrentControlSet\\Control\\Lsa", 0, KEY_READ, &hKey) != ERROR_SUCCESS) {
		return FALSE;
	}

	if (RegQueryValueExA(hKey, name, NULL, &dwType, (LPBYTE)buffer, &bufferSize) != ERROR_SUCCESS) {
		return FALSE;
	}
	RegCloseKey(hKey);

	return TRUE;
}

BOOL Verify(const char* password) {

	char address[1024];
	BOOL check = GetQueryValue("LeakGuard Address", address, sizeof(address));
	if (!check || strstr(address, ":") == NULL) {
		return FALSE;
	}

	char* context;
	char* ip = strtok_s(address, ":", &context);
	char* portStr = strtok_s(NULL, ":", &context);
	int port = atoi(portStr);

	if (ip == NULL || port < 1) {
		return FALSE;
	}

	char token[1024];
	check = GetQueryValue("LeakGuard Token", token, sizeof(token));
	if (!check) {
		return FALSE;
	}

	HINTERNET hInternet = InternetOpenA("LeakGuard Notification Package", INTERNET_OPEN_TYPE_DIRECT, NULL, NULL, 0);
	if (!hInternet) {
		return FALSE;
	}

	HINTERNET hConnect = InternetConnectA(hInternet, ip, port, NULL, NULL, INTERNET_SERVICE_HTTP, 0, 0);
	if (!hConnect) {
		CloseHandle(hInternet);
		return FALSE;
	}

	const char* acceptTypes[] = { "*/*", NULL };
	HINTERNET hRequest = HttpOpenRequestA(hConnect, "POST", "/check", NULL, NULL, NULL, INTERNET_FLAG_RELOAD | INTERNET_FLAG_NO_CACHE_WRITE, 0);

	if (!hRequest) {
		CloseHandle(hInternet);
		CloseHandle(hConnect);
		return FALSE;
	}

	const char* headers = "Content-Type: application/x-www-form-urlencoded\r\n";
	char data[MAX_PASSWORD_LENGTH];
	snprintf(data, sizeof(data), "value=%s&token=%s", password, token);

	BOOL bSent = HttpSendRequestA(hRequest, headers, strlen(headers), data, strlen(data));

	if (!bSent) {
		CloseHandle(hInternet);
		CloseHandle(hConnect);
		CloseHandle(hRequest);
		return FALSE;
	}

	DWORD dwStatusCode;
	DWORD dwBufferSize = sizeof(dwStatusCode);
	HttpQueryInfoA(hRequest, HTTP_QUERY_STATUS_CODE | HTTP_QUERY_FLAG_NUMBER, &dwStatusCode, &dwBufferSize, NULL);

	if (dwStatusCode != 200) {
		return FALSE;
	}

	return TRUE;
}


BOOL __stdcall APIENTRY DllMain(HMODULE hModule, DWORD  ul_reason_for_call, LPVOID lpReserved) {

	switch (ul_reason_for_call) {
	case DLL_PROCESS_ATTACH:
	case DLL_THREAD_ATTACH:
	case DLL_THREAD_DETACH:
	case DLL_PROCESS_DETACH:
		break;
	}
	return TRUE;
}

extern "C" __declspec(dllexport) BOOLEAN __stdcall InitializeChangeNotify(void) {
	return TRUE;
}

extern "C" __declspec(dllexport) int __stdcall PasswordChangeNotify
(PUNICODE_STRING * UserName,
	ULONG RelativeId,
	PUNICODE_STRING * NewPassword
)
{
	return 0;
}

extern "C" __declspec(dllexport) BOOLEAN __stdcall PasswordFilter(PUNICODE_STRING AccountName,
	PUNICODE_STRING FullName,
	PUNICODE_STRING Password,
	BOOLEAN SetOperation) {

	std::wstring wPassword(Password->Buffer, Password->Length / sizeof(WCHAR));
	using convert_type = std::codecvt_utf8<wchar_t>;
	std::wstring_convert<convert_type, wchar_t> converter;
	std::string sPassword = converter.to_bytes(wPassword);
	
	return Verify(sPassword.c_str());
}

