package org.blockchain;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.hyperledger.fabric.shim.ChaincodeBase;
import org.hyperledger.fabric.shim.ChaincodeStub;
import com.google.gson.Gson;

public class LCCreateChainCode extends ChaincodeBase{
	
	private static Log logger = LogFactory.getLog(LCCreateChainCode.class);
	private final String LOGHEAD = "[ LC_CRT_TEST ] ---> ";

	@Override
	public Response init(ChaincodeStub stub) {
		try {
			ldebug("LC create test start!!!");
			if(stub.getArgs().size() != 2){
				return newErrorResponse("Incorrect number parameter!");
			}
			List<String> argList = stub.getParameters();
			String json = argList.get(0);
			Gson gson = new Gson();
			Map<String,String> map = gson.fromJson(json, HashMap.class);
			ldebug(json);
			if(null != map){
				for (Map.Entry<String, String> entry : map.entrySet()) {
					stub.putStringState(entry.getKey(), entry.getValue());
				}
			}
		} catch (Throwable e) {
			return newErrorResponse(e);
		}
		return newSuccessResponse();
	}

	@Override
	public Response invoke(ChaincodeStub stub) {
		try {
			String fun = stub.getFunction();
			List<String> args = stub.getParameters();
			switch(fun){
			case "invoke":
				return doInvoke(stub, args);
			case "query":
				return doQuery(stub, args);
			case "delete":
				return doDelete(stub, args);
			default:
				return newErrorResponse("Unknown function: " + fun);
			}
		} catch (Throwable e) {
			return newErrorResponse(e);
		}
	}
	
	private Response doQuery(ChaincodeStub stub,List<String> argList){
		Map<String,String> data = new HashMap<>();
		if(argList != null && !argList.isEmpty()){
			for (String key : argList) {
				String str = stub.getStringState(key);
				if(null != str){
					data.put(key, str);
				}
			}
		}
		return newSuccessResponse(new Gson().toJson(data));
	}
	
	/**
	 * @Title: doInvoke 
	 * @Description: TODO(修改串) 
	 * @param @param stub
	 * @param @param argList
	 * @param @return    设定参数 
	 * @return Response    返回类型 
	 * @throws
	 */
	private Response doInvoke(ChaincodeStub stub,List<String> argList){
		if(null == argList || argList.size() != 2) throw new IllegalArgumentException("Incorrect number of arguments. Expecting: doInvoke(key, value)");
		String key = argList.get(0);
		String value = argList.get(1);
		stub.putStringState(key, value);
		ldebug("UPDATE" + key + "=" + value);
		return newSuccessResponse();
	}
	
	private Response doDelete(ChaincodeStub stub,List<String> argList){
		if(argList != null && !argList.isEmpty()){
			for (String key : argList) {
				stub.delState(key);
				ldebug("DELETE" + key);
			}
		}
		return newSuccessResponse();
	}
	
	
	/**
	 * @Title: ldebug 
	 * @Description: TODO(打印日志) 
	 * @param @param message    设定参数 
	 * @return void    返回类型 
	 * @throws
	 */
	private void ldebug(String message){
		logger.debug(LOGHEAD + message);
		System.out.println(LOGHEAD + message);
	}
	
	public static void main(String[] args) {
		new LCCreateChainCode().start(args);
	}
}
