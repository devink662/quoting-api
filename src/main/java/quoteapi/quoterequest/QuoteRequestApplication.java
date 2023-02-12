package quoteapi.quoterequest;

import org.flowable.engine.*;
import org.flowable.engine.history.HistoricActivityInstance;
import org.flowable.engine.impl.cfg.StandaloneProcessEngineConfiguration;
import org.flowable.engine.repository.Deployment;
import org.flowable.engine.repository.ProcessDefinition;
import org.flowable.engine.runtime.ProcessInstance;
import org.flowable.task.api.Task;

import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
public class QuoteRequestApplication {


	public static void main(String[] args) {

		// Instantiate a ProcessEngine instance
		ProcessEngineConfiguration cfg = InitiateProcessEngineInstance();

		ProcessEngine processEngine = cfg.buildProcessEngine();


		TaskService taskService = processEngine.getTaskService();

		DeployProcessDefinition(processEngine);

		ProcessInstance processInstance = InitiateProcessInstance(processEngine);

		List<Task> tasks = RetrieveTasks(taskService);

		CompleteTasks(tasks, taskService);

		HistoryService historyService = processEngine.getHistoryService();
		List<HistoricActivityInstance> activities =
				historyService.createHistoricActivityInstanceQuery()
						.processInstanceId(processInstance.getId())
						.finished()
						.orderByHistoricActivityInstanceEndTime().asc()
						.list();

		for (HistoricActivityInstance activity : activities) {
			System.out.println(activity.getActivityId() + " took "
					+ activity.getDurationInMillis() + " milliseconds");
		}
	}

	public static ProcessEngineConfiguration InitiateProcessEngineInstance() {
		// Instantiate a ProcessEngine instance. This is a thread-safe object that you typically have to instantiate
		// A ProcessEngine is created from a ProcessEngineConfiguration instance, which allows you to configure and
		//      tweak the settings for the process engine
		// The minimum configuration a ProcessEngineConfiguration needs is a JDBC connection to a database:
		return new StandaloneProcessEngineConfiguration()
				.setJdbcUrl("jdbc:h2:mem:flowable;DB_CLOSE_DELAY=-1")
				.setJdbcUsername("sa")
				.setJdbcPassword("")
				.setJdbcDriver("org.h2.Driver")
				.setDatabaseSchemaUpdate(ProcessEngineConfiguration.DB_SCHEMA_UPDATE_TRUE);
	}

	public static void DeployProcessDefinition(ProcessEngine processEngine) {
		//  Deploy a process definition to the Flowable engine so that:
		//  1. the process engine will store the XML file in the database, so it can be retrieved whenever needed
		//  2. the process definition is parsed to an internal, executable object model, so that process instances can be started from it.
		RepositoryService repositoryService = processEngine.getRepositoryService();
		Deployment deployment = repositoryService.createDeployment()
				.addClasspathResource("holiday-request.bpmn20.xml")
				.deploy();

		// verify that the process definition is known to the engine
		// Query it through the API by creating a new ProcessDefinitionQuery
		ProcessDefinition processDefinition = repositoryService.createProcessDefinitionQuery()
				.deploymentId(deployment.getId())
				.singleResult();

		System.out.println("Found process definition : " + processDefinition.getName());
	}

	public static ProcessInstance InitiateProcessInstance(ProcessEngine processEngine) {
		// Set Process Instance Variables

		String employee;
		String nrOfHolidays;
		String description;

		Scanner scanner = new Scanner(System.in);

		System.out.println("Who are you?");
		employee = scanner.nextLine();

		System.out.println("How many holidays do you want to request?");
		nrOfHolidays = String.valueOf(Integer.valueOf(scanner.nextLine()));

		System.out.println("Why do you need them?");
		description = scanner.nextLine();

		Map<String, Object> variables = new HashMap<>();
		variables.put("employee", employee);
		variables.put("nrOfHolidays", nrOfHolidays);
		variables.put("description", description);


		// start a process instance through the RuntimeService. The collected data is passed as a java.util.Map instance,
		// where the key is the identifier that will be used to retrieve the variables later on. The process instance is
		// started using a key. This key matches the id attribute that is set in the BPMN 2.0 XML file,
		// in this case holidayRequest.

		RuntimeService runtimeService = processEngine.getRuntimeService();

		ProcessInstance processInstance =
				runtimeService.startProcessInstanceByKey("holidayRequest", variables);

		return processInstance;
	}

	public static List<Task> RetrieveTasks(TaskService taskService) {

		List<Task> tasks = taskService.createTaskQuery().taskCandidateGroup("managers").list();
		System.out.println("You have " + tasks.size() + " tasks:");
		for (int i = 0; i < tasks.size(); i++) {
			System.out.println((i + 1) + ") " + tasks.get(i).getName());
		}
		return tasks;
	}

	public static void CompleteTasks(List<Task> tasks, TaskService taskService) {
		Scanner scanner = new Scanner(System.in);
		System.out.println("Which task would you like to complete?");
		int taskIndex = Integer.valueOf(scanner.nextLine());
		Task task = tasks.get(taskIndex - 1);
		Map<String, Object> processVariables = taskService.getVariables(task.getId());
		System.out.println(processVariables.get("employee") + " wants " +
				processVariables.get("nrOfHolidays") + " of holidays. Do you approve this?");

		boolean approved = scanner.nextLine().equalsIgnoreCase("y");
		Map<String, Object> variables = new HashMap<>();
		variables.put("approved", approved);
		taskService.complete(task.getId(), variables);
	}
}
