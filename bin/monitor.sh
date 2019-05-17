#! /bin/sh  

proc_name="mj"  

proc_num()  
{ 
  num=`ps -ef | grep $proc_name | grep -v grep | wc -l` 
  return $num 
} 

proc_num 
number=$?  

echo $number

if [ $number -eq 0 ] 
then 
  cd /home/nano/gamebackend/bin; ./run.sh
fi 
