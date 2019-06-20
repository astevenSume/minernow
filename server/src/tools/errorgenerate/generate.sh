#!/bin/bash
bin/errorgenerate --config src/tools/errorgenerate/otc.json --out-dir src/otc_error/
bin/errorgenerate --config src/tools/errorgenerate/admin.json --out-dir src/admin/controllers/errcode/
