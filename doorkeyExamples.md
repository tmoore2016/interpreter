# Doorkey Examples
<br/>
*// = comments*<br/>
**Bold = Result**<br/>
<br/>
*// Assign variables, create a function, call the function with parameters*<br/>
let five = 5;<br/>
<br/>
let ten = 10;<br/>
<br/>
let add = fn(x, y) {<br/>
  x + y;<br/>
};<br/>
<br/>
let result = add(five, ten);<br/>
<br/>
result<br/>
**15**<br/>
<br/>
*// Create and insert elements into an array*<br/>
let arr = [1,2,3];<br/>
arr<br/>
**[1, 2, 3]**<br/>
<br/>
*// Call builtin len() to count elements in array*<br/>
len(arr)<br/>
**3**<br/>
<br/>
*// Call the first element in an array*<br/>
first(arr)<br/>
**1**<br/>
<br/>
*// Call the last element in an array*<br/>
last(arr)<br/>
**3**<br/>
<br/>
*// Call all elements except the first, in an array*<br/>
tail(arr)<br/>
**[2,3]**<br/>
<br/>
*// Create a function*<br/>
let square = fn(x){x*x};<br/>
<br/>
*// Call the second array element as a parameter within a function*
square(arr[1])<br/>
**4**<br/>
<br/>
*// Call the builtin push() function to add an element to the array.*
let arr = push(arr, "four");<br/>
arr<br/>
**[1, 2, 3, four]**<br/>
<br/>
*// Use operators on array elements*<br/>
arr[1] * arr[2]<br/>
**6**<br/>
<br/>
*// Concatenate string element inside array*<br/>
arr[3] + "th"<br/>
**fourth**<br/>
<br/>
*// Concatenate string inside array and push into the array*<br/>
let arr = push(arr, arr[5] + "th");<br/>
arr<br/>
**[1,2,3,four,fourth]**<br/>
<br/>
*// Create a hash table*<br/>
let books = [{"title": "The Sea-Wolf", "authorFirstName": "Jack", "authorLastName": "London", "age": 115}, {"title": "Fear and Loathing in Las Vegas", "authorFirstName": "Hunter", "authorLastName": "Thompson", "age": 48}];<br/>
<br/>
books<br/>
**[{title: The Sea-Wolf, authorFirstName: Jack, authorLastName: London, age: 115}, {title: Fear and Loathing in Las Vegas, authorFirstName: Hunter, authorLastName: Thompson, age: 48}]**<br/>
<br/>
books[0]["title"]<br/>
**The Sea-Wolf**<br/>
<br/>
*// Create a mapping function*<br/>
let a = [1,3,5,7];<br/>
<br/>
let square = fn(x) {x*x};<br/>
<br/>
let map = fn(newArr, f) {let iter = fn(newArr, accumulated) { if (len(newArr) == 0) {accumulated} else{iter(tail(newArr),push(accumulated, f(first(newArr))));}}; iter(newArr, []);};<br/>
<br/>
map(a,square);<br/>
**[1, 9, 25, 49]**<br/>
