import bs4

def main():
   # open the guide html exported from Acrobat
   guide = open('RefGuide.html')
   # generate a Beautiful Soup object to operate on
   # I'm using lxml as the parser, 'pip install lxml' if you don't have it installed
   soup = bs4.BeautifulSoup(guide, 'lxml')

   # generate a list of all the headings for each entry in the guide, this is the most reliable pattern in the file
   # we'll use this to walk siblings to get the rest of the entry
   headings = soup.find_all('h1')
   # there should be 675 entries, some of these will just be numbers because junky formatting
   # we need to strip these out. There will still be some headings with numbers at the end, but that's ok
   # removing only number h1 entries allows us to use H1s as the dividing line for each section
   cleaned_headings = []
   for heading in headings:
      if heading.getText().isdigit() == False:
         cleaned_headings.append(heading)
   
   # Now that we have "cleaned" H1 entries, we can group all the still nosiy data between them
   # these are actual entry groupings, although they'll still need work after this
   # this should pull in all but 7 of the entries that have malformed lists in them, becuase of course they do
   # I'll add the missing 7 by hand and call it good enough. 
   entries = []
   for heading in cleaned_headings:
      this_entry = [heading]
      current = heading.find_next_sibling()
      while True:
         try:
            if current.name != 'h1':
               this_entry.append(current)
               current = current.find_next_sibling()
            elif current.name == 'h1' and current.getText().isdigit() == True:
               this_entry.append(current)
               current = current.find_next_sibling()
            else:
               entries.append(this_entry)
               break
         except AttributeError as error:
            print(" It broke at " + str(heading))
            break
   # Strip out all the actual text from the rest of the cruft now that tags are no longer useful
   raw_guide_text = []
   for entry in entries:
      this_entry = []
      for line in entry:
         if line.getText().isprintable():
            this_entry.append(line.getText())
      raw_guide_text.append(this_entry)

   # Things we know for sure about the data in raw_guide_text at this point:
   # the entries are not all the same length, smallest one is 10 items, longest is 30
   # the first item in each entry will always be the title of the resource
   # the second item will be an address with just a few exceptions
   # the the third item is sometimes the address, but most often a number or the phone column label
   # phone numbers should always appear in the order "Crisis line: Business line: Fax: TDD:"
   # however, there isn't any indication which numbers are missing if there are < 4 in the set
   # the desctiption ALWAYS follows the 'Description:' entry
   # however, 14 of the entries don't have a 'Description:' field, but appear to have descriptions
   # some of the time the services offered follows the 'Services:' entry
   # however, 14 of the entries don't have a 'Services:' entry
   # and in 55 entries 'Wheelchair Accessible: Languages' followes the 'Services:' entry
   # Wheelchair Accessible, Languages, Populations served, and Categories seem to follow the same pattern as services and description
   # the website URL will always be in an entry that starts with 'Website'
   # the email address will always be in an entry that starts with 'Email'
   
   new_guide = []
   for entry in raw_guide_text:
      new_entry = {}
      new_entry["name"] = entry[0]
      # address may not be correct. This will just need human eyes on it
      new_entry["address"] = entry[1]
      # if no description exists in the guide, we'll mark it "needs description"
      try:
         f = entry.index('Description:')
         new_entry["description"] = entry[f+1]
      except ValueError as error:
         new_entry["description"] = "This entry needs a description"
      # if no services exists or is the weird junk data we'll mark it "needs services"
      # there is also a weird chance we might have the whole wheelchair entry in here,
      # if so, we'll just add it in right now and check for it in the wheelchar phase
      try:
         f = entry.index('Services:')
         if entry[f+1] == "Services: Wheelchair Accessible: Yes":
            new_entry["wheelchair"] = True
            new_entry["services"] = "This entry needs a services list"
         # these are the junk entries that might exist in this field
         elif entry[f+1] == "Wheelchair Accessible:" 
            new_entry["services"] = "This entry needs a services list"
         # this one tells us there isn't a languges list either
         elif entry[f+1] == "Wheelchair Accessible: Languages"
            new_entry["languges"] = "This entry nees a languages list"
            new_entry["services"] = "This entry needs a services list"
         else:
            new_entry["services"] = entry[f+1]
      except ValueError as error:
          new_entry["services"] = "This entry needs a services list"





if __name__ == '__main__':
   main()