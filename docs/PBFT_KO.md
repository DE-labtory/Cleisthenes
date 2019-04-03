<h2>PBFT ����<h2/>

<h4>�л� �ý��� ����</h4> 

<ul>
<li> Fail-stop</li>
	
������ ���� �浹 ���������� �浹 �� �ٸ� ������ �浹������ �߰� ����,���� return ������ �ʴ´�.	
<li> Byzantine fault</li>
	
Fail-stop ó�� ������ ���� �浹 ����������, �浹���� �� �ٸ� ������ �߰��ϱ� ���� �ʰ�, 
  �߸��� ������ �浹�� ������ return �Ѵ�.
</ul>

<h4>Byzantine ������?</h4>

<p align="center"><img src="../images/pbft_byzantine.PNG" width="400px" height="200px"></p>



<p>��ü ��ũ��ũ�� ����� ���� 3����� ������ �غ����� ��, �� �� commander ������ �ϴ� ��尡 1�� �����ϰ� �ȴ�. 
Commander ���� ��ǥ�� ������ ������ �����Ͽ� ������ Lieutenant ������ ��ǥ�� �� �� �ְ� �� ��, �� ����� �ݿ��Ͽ� �ൿ�� �����ϴ� ������ �Ѵ�. 
�� 3���� ��尡 ��� �ùٸ� ������ ���ο��� �����Ѵٸ� 3���� ���� consensus�� ������ �� �ְ� �ȴ�. 
�׷��� 3���� ��� �� �ܺ��� �������� ���� ������ ������ �����ϰų� �������� �����ϴ� ���� �����ϴ� ������ 
�����ϴ� ��尡 ���� �� �ִ�. (Byzantine ���� ���� ����, 1���ϼ���, �� �̻��� ���� �ִ�.) �� �� �� ��带 
Byzantine ����� �Ѵ�. �� Byzantine ���� ���� ������ �߻��ϰ� �ȴ�. </p>

<p>���� Byzantine ��尡 1����� ������ 2���� ��� �� 1���� ������ ������ �ùٸ� ������ ���� �����ް� �Ǿ� ȥ���� �Ͼ��
 �ȴ�. �� �� ��� ���� �ùٸ� ���������� �� ����� ���� ������ ������ �߻��Ѵ�. 


���� ���ø� Ȯ���Ͽ� PBFT�� ���� �˾ƺ����� ����. </p>


<h4>PBFT</h4>

<p>PBFT�� �л�ý����� Byzantine ��尡 ���� �� �� �ִ� �񵿱� �ý����� �� �ش� �л�ý��ۿ� 
 ������ ��� ��尡 ���������� ���Ǹ� �̷� �� �ֵ��� ���ߵ� ���� �˰������̴�. ������ BFT ���� �˰������� ����� 
 ��Ʈ��ũ������ ���ǰ� �����ߴ� ������ �����Ͽ� Byzantine ��尡 �ִ� �񵿱� ��ũ��ũ������ ������ �� �ְ� �Ͽ���.</p>
 

<p>PBFT�� ��ü ��ũ��ũ ����� ���� 3*F+1 �϶�, F���� Byzantine ������ ����� �� �ִٴ� ������ �ϰ� �ִ�. 
 ������ ����Ѵٴ� ���� Byzantine ������ ���� �ʰ� ����� �� ���� return�ϴ� ���� �ǹ��Ѵ�. PBFT�� client�� 
 ù��° ��忡�� request�� ������ �� ���۵ȴ�. �� ù��° ���� client�� ��û�� ��� transaction�� ��� ������ 
 �����. �� ���� 3�ܰ踦 ���� consensus�� �̸��� �ȴ�.</p>
<p align="center"><img src="../images/pbft_3steps.png" width="400px" height="200px"></p>



<h4>PBFT �ֿ� 3�ܰ�</h4>

<ul>
<li> Pre prepare</li>
	
Client�� ��û�� ���� ù��° ��尡 ��Ʈ��ũ�� �ִ� ��� ��忡�� ������ �����ϰ�(pre prepare �޼����� ������),

��� ������ �� ������ �ް� �ȴ�. Pre prepared �޼������� ������ �޼����� ��ȿ���� ������ �� �ִ� signiture�� �ٸ� meta data, sequence number�� ��� �ִ�.   
	
<li> Prepare</li>
	
��� ������ ������ �ް� �Ǹ� ������ ���� ��尡 ������ ��� ���鿡�� �ڽ��� ������ �޾Ҵٴ� prepare �޼����� ������ �ȴ�. �� ���� �ٸ� ������ ������ �޾Ҵ��� ���θ� �����Ͽ� �� ���� 2/3 �̻��� �� ������ �����Ѵ�. 

 (Prepare �޼����� �� �޼����� ���� ��尡 ���� ��ȿ�ϴٴ� ���� �˷��ش�.)   
	
<li> Commit</li>
	
�� �� ���� ������� ��� ��忡 �����Ѵ�. �� ���� �ٸ� ������ ������ ���� ��ȿ�� ���� ������� �����ϸ�, ��ü�� 2/3 �ʰ��Ͽ� ������ ������� ������ ��� �ش� ������� ������ �ν��Ѵ�. 
	
-���� ������� 2/3�� �ʰ��� ��尡 ������ ��ȿ�ϴٴ� ������� ������ ��, ������ �ڽ��� ����ü�ο� �߰��Ѵ�.
	
-1/3 �̻��� ��尡 ������ ��ȿ���� �ʴٴ� ������� ������ ��, ������ ����ü�ο� �߰����� �ʴ´�.
	
�� ������ ������ ���� ���¸� client���� reply ���ش�.
	</ul>
	

