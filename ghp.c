#include <stdio.h>
#include <stdlib.h>
#include "lib/ph7.h"
#include "ghp.h"

extern int gph_output_consumer(const void *pOutput, unsigned int nOutputLen, void *pUserData);
extern void gph_error_consumer(const char *zMsg);

int gph_eval(const char *zSource, const char **args, int argc, const gph_foreign_func *ffs, int ffLen) {
	ph7 *pEngine;
	ph7_vm *pVm;
	int rc;
	int i;
	rc = ph7_init(&pEngine);
	if(rc != PH7_OK) {
		gph_error_consumer("Error while allocating");
        goto Fatal;
	}

	ph7_config(pEngine, PH7_CONFIG_ERR_OUTPUT, gph_output_consumer, 0);

	rc = ph7_compile_v2( pEngine, zSource, -1, &pVm, 0);
	if(rc != PH7_OK) {
		if(rc == PH7_COMPILE_ERR) {
			const char *zErrLog;
			int nLen;
			ph7_config(pEngine, PH7_CONFIG_ERR_LOG, &zErrLog, &nLen);
			if(nLen > 0){
				gph_error_consumer(zErrLog);
			}
		}
		gph_error_consumer("Compile error");
        goto Fatal;
	}

	rc = ph7_vm_config(pVm, PH7_VM_CONFIG_OUTPUT, gph_output_consumer, 0);
	if(rc != PH7_OK) {
		gph_error_consumer("Error installing output consumer callback");
		goto Fatal;
	}

    for(i = 0; i < ffLen;  i++) {
        rc = ph7_create_function(pVm, ffs[i].zName, ffs[i].xProc, 0);
    	if(rc != PH7_OK){
    		gph_error_consumer("Error registering foreign functions");
            goto Fatal;
    	}
    }

    // 开启运行时错误报告
	ph7_vm_config(pVm, PH7_VM_CONFIG_ERR_REPORT);

    // 填充 $argv
    for(i = 0; i < argc;  i++) {
        ph7_vm_config(pVm, PH7_VM_CONFIG_ARGV_ENTRY, args[i]);
    }

	ph7_vm_exec(pVm, 0);
	ph7_vm_release(pVm);
	ph7_release(pEngine);
	return 0;

	Fatal:
	ph7_lib_shutdown();
    return -1;
}