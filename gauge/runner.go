package gauge

import (
	"fmt"
	"net"
	"os"
	"reflect"

	"regexp"

	c "github.com/manuviswam/gauge-go/constants"
	m "github.com/manuviswam/gauge-go/gauge_messages"
	mp "github.com/manuviswam/gauge-go/messageprocessors"
	mu "github.com/manuviswam/gauge-go/messageutil"
	t "github.com/manuviswam/gauge-go/testsuit"
)

var context *t.GaugeContext
var processors mp.ProcessorDictionary

func init() {
	context = &t.GaugeContext{
		Steps:         make([]t.Step, 0),
		Hooks:         make([]t.Hook, 0),
		SuiteStore:    make(map[string]interface{}),
		SpecStore:     make(map[string]interface{}),
		ScenarioStore: make(map[string]interface{}),
	}

	processors = mp.ProcessorDictionary{}
	processors[*m.Message_StepNamesRequest.Enum()] = &mp.StepNamesRequestProcessor{}
	processors[*m.Message_StepValidateRequest.Enum()] = &mp.StepValidateRequestProcessor{}
	processors[*m.Message_SuiteDataStoreInit.Enum()] = &mp.SuiteDataStoreInitRequestProcessor{}
	processors[*m.Message_ExecutionStarting.Enum()] = &mp.ExecutionStartingRequestProcessor{}
	processors[*m.Message_SpecExecutionStarting.Enum()] = &mp.SpecExecutionStartingRequestProcessor{}
	processors[*m.Message_ScenarioExecutionStarting.Enum()] = &mp.ScenarioExecutionStartingRequestProcessor{}
	processors[*m.Message_StepExecutionStarting.Enum()] = &mp.StepExecutionStartingRequestProcessor{}
	processors[*m.Message_ExecuteStep.Enum()] = &mp.ExecuteStepProcessor{}
	processors[*m.Message_ExecutionEnding.Enum()] = &mp.ExecutionEndingProcessor{}
	processors[*m.Message_StepExecutionEnding.Enum()] = &mp.StepExecutionEndingProcessor{}
	processors[*m.Message_ScenarioExecutionEnding.Enum()] = &mp.ScenarioExecutionEndingProcessor{}
	processors[*m.Message_SpecExecutionEnding.Enum()] = &mp.SpecExecutionEndingProcessor{}
	processors[*m.Message_SpecDataStoreInit.Enum()] = &mp.SpecDataStoreInitProcessor{}
	processors[*m.Message_ScenarioDataStoreInit.Enum()] = &mp.ScenarioDataStoreInitProcessor{}
}

func BeforeSuite(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESUITE,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func AfterSuite(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSUITE,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func BeforeSpec(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESPEC,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func AfterSpec(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSPEC,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func BeforeScenario(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESCENARIO,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func AfterScenario(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSCENARIO,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func BeforeStep(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.BEFORESTEP,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func AfterStep(fn func(), tags []string, op t.Operator) bool {
	hook := t.Hook{
		Type:     t.AFTERSTEP,
		Impl:     fn,
		Tags:     tags,
		Operator: op,
	}

	context.Hooks = append(context.Hooks, hook)
	return true
}

func GetSuiteStore() map[string]interface{} {
	return context.SuiteStore

}

func GetSpecStore() map[string]interface{} {
	return context.SpecStore
}

func GetScenarioStore() map[string]interface{} {
	return context.ScenarioStore
}

func Run() {
	fmt.Println("We have got ", len(context.Steps), " step implementations") // remove

	var gaugePort = os.Getenv(c.GaugePortVariable)

	fmt.Println("Connecting port:", gaugePort) // remove
	conn, err := net.Dial("tcp", net.JoinHostPort("127.0.0.1", gaugePort))
	defer conn.Close()
	if err != nil {
		fmt.Println("dial error:", err)
		return
	}
	for {
		msg, err := mu.ReadMessage(conn)
		if err != nil {
			fmt.Println("Error reading message : ", err)
			return
		}
		if *msg.MessageType.Enum() == *m.Message_KillProcessRequest.Enum() {
			return
		}

		processor := processors[*msg.MessageType.Enum()]

		if processor == nil {
			fmt.Println("Unable to find processor for message type : ", msg.MessageType)
			return
		}
		msgToSend := processor.Process(msg, context)

		mu.WriteGaugeMessage(msgToSend, conn)
	}
}

func Step(stepDesc string, impl interface{}) bool {
	desc, noOfArgs := parseDesc(stepDesc)
	implType := reflect.TypeOf(impl)

	if reflect.ValueOf(impl).Kind() != reflect.Func {
		//TODO decide whether to ignore or fail test
		fmt.Printf("Expected a function implementation for '%s' but got type '%s' - Ignoring test\n", stepDesc, implType.String())
		return false
	}

	//TODO validate not just the number of arguments but method signature
	if implType.NumIn() != noOfArgs {
		//TODO decide whether to ignore or fail test
		fmt.Printf("Mismatch in number of arguments in implementation of '%s' expected : %d, actual : %d - Ignoring test\n", desc, noOfArgs, implType.NumIn())
		return false
	}
	step := t.Step{
		Description: desc,
		Impl:        impl,
	}
	context.Steps = append(context.Steps, step)
	return true
}

func parseDesc(desc string) (string, int) {
	re := regexp.MustCompile("<(.*?)>")
	return re.ReplaceAllLiteralString(desc, "{}"), len(re.FindAllString(desc, -1))
}
