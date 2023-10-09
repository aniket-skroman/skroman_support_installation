# skroman support for installation department 

 ```migrate
 install migrate first for DB migration
   1. Windows:
        a. Open the Windows power shell on your PC, and type the command
            
            **irm get.scoop.sh | iex**

        b. then install the migrate

            **scoop install migrate**  

   2. Ubuntu:
        a. Let us setup the repository to install the migrate package.
            
            **curl -s https://packagecloud.io/install/repositories/golang-migrate/migrate/script.deb.sh | sudo bash **
        b. Update the system by executing the following command

            **sudo apt-get update**

        c. Now, it’s time to set up golang migrate.

            **sudo apt-get install migrate**   

 ``` 
  